package main

import (
	"fmt"
	"time"
)

type Ball struct {
	step int
}

func main() {
	ch := make(chan *Ball)
	go player("Ahmad", ch)
	go player("Haris", ch)
	ch <- &Ball{}
	time.Sleep(1 * time.Second)
	<- ch
}

func player(name string, ch chan *Ball) {
	for {
		ball := <-ch
		fmt.Println(name, "playing the ball", ball.step, "now")
		ball.step++
		time.Sleep(100 * time.Millisecond)
		ch <- ball
	}
}