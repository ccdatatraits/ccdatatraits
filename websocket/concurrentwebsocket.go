package main

import (
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"github.com/gorilla/websocket"
)

const (
	rootHTMLLocation = "html/index.html"
)

type concurrentSocketHandler struct {
	websocketUpgrader websocket.Upgrader
	msgChan chan string
	done chan bool
}
/*func (socketHandler *concurrentSocketHandler) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := socketHandler.websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		if string(msg) == "ping" {
			fmt.Println("ping")
			err = conn.WriteMessage(msgType, []byte("pong"))
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			conn.Close()
			fmt.Println(string(msg))
			return
		}
	}
}*/
func (socketHandler *concurrentSocketHandler) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := socketHandler.websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		iMsg := string(msg)
		if iMsg == "ping" {
			socketHandler.msgChan <- iMsg
			err = conn.WriteMessage(msgType, []byte(<-socketHandler.msgChan))
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			socketHandler.done <- true
			conn.Close()
			fmt.Println(string(msg))
			return
		}
	}
}

func NewConcurrentSocketHandler() *concurrentSocketHandler {
	newSocketHandler := new(concurrentSocketHandler)
	newSocketHandler.websocketUpgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	newSocketHandler.msgChan = make(chan string)
	newSocketHandler.done = make(chan bool)
	return newSocketHandler
}
func (socketHandler *concurrentSocketHandler) GetWebsocketHandler() func(http.ResponseWriter, *http.Request) {
	return socketHandler.handler;
}
func main() {
	http.HandleFunc("/", rootHandler)
	socketHandler := NewConcurrentSocketHandler()
	go socketRoutine(socketHandler)
	http.HandleFunc("/websocket", socketHandler.GetWebsocketHandler())
	log.Fatal(http.ListenAndServe(":3000", nil))
}
func socketRoutine(socketHandler *concurrentSocketHandler) {
	defer func() {close(socketHandler.msgChan); close(socketHandler.done)}()
	for {
		select {
		case iMsg := <-socketHandler.msgChan:
			fmt.Println(iMsg)
			switch iMsg {
			case "ping":
				socketHandler.msgChan <- "pong"
			case "pong":
				socketHandler.msgChan <- "ping"
			default:
				socketHandler.msgChan <- iMsg
			}
		case <-socketHandler.done:
			return
		}
	}
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	indexFile, err := os.Open(rootHTMLLocation)
	if err != nil {
		fmt.Println(err)
	}
	index, err := ioutil.ReadAll(indexFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(index))
}
