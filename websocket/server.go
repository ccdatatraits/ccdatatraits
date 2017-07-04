package main

import (
	"fmt"
	"strings"
	"net/http"
	"os"
	"io/ioutil"
	"log"
	"time"
	"github.com/gorilla/websocket"
)

const (
	DEFAULT_LISTEN_ADDR        = ":3000"
	DEFAULT_ROOT_HTML_LOCATION = "html/index.html"
)

type ServerConfig struct {
	ListenAddr       string
	RootHTMLLocation string
}
type ServerService interface {
	Configure(ServerConfig) ServerService
	NewServer() Server
}
type Server interface {
	Start() error
	GetCommChannel() (chan string)
	Stop() error
}

type serverServiceImpl struct {
	savedConfiguration ServerConfig
}

func (service *serverServiceImpl) Configure(config ServerConfig) ServerService {
	savedServerListenAddr := strings.Trim(config.ListenAddr, " ")
	if "" == savedServerListenAddr {
		service.savedConfiguration.ListenAddr = DEFAULT_LISTEN_ADDR
	} else {
		service.savedConfiguration.ListenAddr = savedServerListenAddr
	}
	savedServerRootHTMLLocation := strings.Trim(config.RootHTMLLocation, " ")
	if "" == savedServerRootHTMLLocation {
		service.savedConfiguration.RootHTMLLocation = DEFAULT_ROOT_HTML_LOCATION
	} else {
		service.savedConfiguration.RootHTMLLocation = savedServerRootHTMLLocation
	}
	return service
}
func (service *serverServiceImpl) NewServer() Server {
	newServerImpl := new(serverImpl)
	newServerImpl.configuration = service.savedConfiguration
	return newServerImpl
}

type httpResWriterReqPointer struct {
	w http.ResponseWriter
	r *http.Request
}

type serverImpl struct {
	configuration ServerConfig
	serverChannel chan string
	serverStopped bool
	stopChannel   chan struct{}
	conn          httpResWriterReqPointer
	done          chan bool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func (server *serverImpl) Start() error {
	fmt.Println("Configured for address:", server.configuration.ListenAddr)
	fmt.Println("Root HTML location:", server.configuration.RootHTMLLocation)
	fmt.Println("Running server now")
	fmt.Println("")
	http.HandleFunc("/", server.rootHandler)
	http.HandleFunc("/websocket", server.socketHandler)
	err := http.ListenAndServe(server.configuration.ListenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
func (server *serverImpl) rootHandler(w http.ResponseWriter, r *http.Request) {
	indexFile, err := os.Open(server.configuration.RootHTMLLocation)
	if err != nil {
		fmt.Println(err)
	}
	index, err := ioutil.ReadAll(indexFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(index))
}
func (server *serverImpl) socketHandler(w http.ResponseWriter, r *http.Request) {
	server.serverChannel = make(chan string)
	server.stopChannel = make(chan struct{})
	server.conn = httpResWriterReqPointer{w, r}
	server.done = make(chan bool)
	go server.serverChannelHandler()
	<-server.done
}
func (server *serverImpl) GetCommChannel() (chan string) {
	return server.serverChannel
}
func (server *serverImpl) serverChannelHandler() {
	const (
		TOTAL = 10
	)
	var OUT_MSGS = [2]string{"ping", "pong"}
	i := 0
	conn, err := upgrader.Upgrade(server.conn.w, server.conn.r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
HERE:
	for {
		select {
		case <-server.stopChannel:
			fmt.Println("Received: STOP! (signal)")
			close(server.serverChannel)
			server.serverStopped = true
			break HERE
		case incoming := <-server.serverChannel:
			fmt.Println("Received:", incoming)
			switch incoming {
			case "PING!":
				server.serverChannel <- "PONG!"
			case "ALIVE?":
				server.serverChannel <- "YES!"
			default:
				server.serverChannel <- incoming
			}
		default:
			msgType, msg, err := conn.ReadMessage()
			incoming_msg := string(msg)
			if err != nil {
				fmt.Println(err)
				return
			}
			if incoming_msg == "ping" || incoming_msg == "pong" {
				fmt.Println(incoming_msg)
				time.Sleep(500 * time.Millisecond)
				err = conn.WriteMessage(msgType, []byte(OUT_MSGS[i%2]))
				if err != nil {
					fmt.Println(err)
					return
				}
			} else {
				conn.Close()
				fmt.Println(incoming_msg)
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
		i++
	}
	fmt.Println("Server Comm Channel stopped!")
	//server.stopChannel <- struct{}{}
	//close(server.stopChannel)
	conn.Close()
	<-server.done
	close(server.done)
}
func (server *serverImpl) Stop() error {
	server.stopChannel <- struct{}{}
	<-server.stopChannel
	return nil
}
func handleMessage(server Server, msg string) string {
	server.GetCommChannel() <- msg
	return "Sent back:" + <-server.GetCommChannel()
}

func main() {
	server := new(serverServiceImpl).Configure(ServerConfig{}).NewServer()
	server.Start()
	//fmt.Println(handleMessage(server, "PING!"))
	//fmt.Println(handleMessage(server, "ALIVE?"))
	//server.Stop()

	fmt.Println("\nENDED!")
}
