package main

import (
	"fmt"
	"bufio"
	"os"
)

func readString(channel <-chan string) {
	filename_read := false
	for {
		str_read := <-channel
		if !filename_read {
			fmt.Println("filename read: " + str_read)
			filename_read = !filename_read
		} else {
			fmt.Println("string read: " + str_read)
		}
	}
}

func main() {
	channel := make(chan string)
	go readString(channel)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter filename:")
	text, _ := reader.ReadString('\n')
	channel <- text
	close(channel)
}