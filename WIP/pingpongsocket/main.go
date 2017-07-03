package main

import (
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"time"
)

const listenAddr = "localhost:3000"

type socket struct {
	io.ReadWriter
	done chan bool
}

func (s socket) Close() error {
	s.done <- true
	return nil
}

func socketHandler(ws *websocket.Conn) {
	s := socket{ws, make(chan bool)}
	go runPingPong(s)
	<- s.done
}
func runPingPong(conn io.ReadWriteCloser) {
	msgByteArray := []byte{}
	const (
		TOTAL = 10
	)
	var OUT_MSGS = [2]string{"ping", "pong"}
	i := 0
	for {
		bytesRead, err := conn.Read(msgByteArray)
		if err != nil {
			fmt.Println("Issue reading:", err)
		}
		fmt.Printf("%v bytes read: (%v)\n", bytesRead, string(msgByteArray))
		time.Sleep(500 * time.Millisecond)
		outMsg := OUT_MSGS[i%2]
		bytesWritten, err := conn.Write([]byte(outMsg))
		if err != nil {
			fmt.Println("Issue writing:", err)
		}
		fmt.Printf("%v bytes written: (%v)\n", bytesWritten, outMsg)
		if i >= TOTAL {
			break
		}
		i++
	}
	conn.Close()
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/websocket", websocket.Handler(socketHandler))
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	indexFile, err := os.Open("WIP/pingpongsocket/html/index.html")
	if err != nil {
		fmt.Println(err)
	}
	index, err := ioutil.ReadAll(indexFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(index))
}