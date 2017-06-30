package main

import (
	"fmt"
	"time"
)

func main() {
	messagebuspasser := make(chan chan int)
	go consumer(messagebuspasser)
	go producer(messagebuspasser)
	time.Sleep(300 * time.Millisecond)
	//panic("show me the stacks")
}

func consumer(messagebuspasser chan<- chan int) {
	messagebus := make(chan int)
	// First pass the message bus onto the passer
	messagebuspasser <- messagebus
	// then wait for all messages (range messagebus)
	for i := range messagebus {
		fmt.Println(i)
	}
}

func producer(messagebuspasser <-chan chan int) {
	// First get the messagebus sent by the consumer
	messagebus := <- messagebuspasser
	// Then push data onto it
	for i := 0; i < 10; i++ {
		//time.Sleep(100 * time.Millisecond)
		messagebus <- i
	}
	close(messagebus)
}
