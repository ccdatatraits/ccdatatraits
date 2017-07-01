package main

import (
	"net"
	"log"
	"io"
)

const listenAddr = "localhost:4000"

func main() {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go io.Copy(c, c)
	}
}
