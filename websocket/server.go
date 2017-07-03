package main

import (
	"fmt"
	"strings"
)

const (
	DEFAULT_LISTENADDR = ":3000"
	DEFAULT_ROOTHTMLLOCATION = "html/index.html"
)

type ServerConfig struct {
	ListenAddr string
	RootHTMLLocation string
}
type ServerService interface {
	Configure(ServerConfig) ServerService
	NewServer() Server
}
type Server interface {
	Start() Server
	GetCommChannel() (chan string)
	HasStopped() bool
	Stop() string
}

type ServerServiceImpl struct {
	savedConfiguration ServerConfig
}
func (service *ServerServiceImpl) Configure(config ServerConfig) ServerService {
	savedServerListenAddr := strings.Trim(config.ListenAddr, " ")
	if "" == savedServerListenAddr {
		service.savedConfiguration.ListenAddr = DEFAULT_LISTENADDR
	} else {
		service.savedConfiguration.ListenAddr = savedServerListenAddr
	}
	savedServerRootHTMLLocation := strings.Trim(config.ListenAddr, " ")
	if "" == savedServerRootHTMLLocation {
		service.savedConfiguration.RootHTMLLocation = DEFAULT_ROOTHTMLLOCATION
	} else {
		service.savedConfiguration.RootHTMLLocation = savedServerRootHTMLLocation
	}
	return service
}
func (service *ServerServiceImpl) NewServer() Server {
	newServerImpl := new(ServerImpl)
	newServerImpl.configuration = service.savedConfiguration
	return newServerImpl
}

type ServerImpl struct {
	configuration ServerConfig
	serverChannel chan string
	serverStopped bool
	stopChannel chan struct{}
}
func (server *ServerImpl) Start() Server {
	fmt.Println("Configured for address:", server.configuration.ListenAddr)
	fmt.Println("Root HTML location:", server.configuration.RootHTMLLocation)
	fmt.Println("Running server now\n")
	server.serverChannel = make(chan string)
	server.stopChannel = make(chan struct{})
	go server.serverChannelHandler()
	return server
}
func (server *ServerImpl) GetCommChannel() (chan string) {
	return server.serverChannel
}
func (server *ServerImpl) serverChannelHandler() {
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
			case "STOP!":
				go server.Stop()
			default:
				server.serverChannel <- incoming
			}
		}
	}
	fmt.Println("Server Comm Channel stopped!")
	server.stopChannel <- struct{}{}
	close(server.stopChannel)
}
func (server *ServerImpl) HasStopped() bool {
	return server.serverStopped
}
func (server *ServerImpl) Stop() string {
	if server.HasStopped() {
		return "Server has already been stopped!"
	}
	server.stopChannel <- struct{}{}
	return "Sent back:STOPPED!"
}
func handleMessage(server Server, msg string) string {
	if server.HasStopped() {
		return "Server has already been stopped!"
	}
	server.GetCommChannel() <- msg
	return "Sent back:" + <-server.GetCommChannel()
}

func main() {
	server := new(ServerServiceImpl).Configure(ServerConfig{}).NewServer()
	serverCommChan := server.Start()
	fmt.Println(handleMessage(serverCommChan, "PING!"))
	fmt.Println(handleMessage(serverCommChan, "ALIVE?"))
	fmt.Println(handleMessage(serverCommChan, "STOP!"))
	fmt.Println(server.Stop())

	fmt.Println("\nENDED!")
}
