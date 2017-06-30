package main

import (
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type CustomIntSplit interface {
	Split(data []byte, atEOF bool) (advance int, token []byte, err error)
	Fetch() int64
}

type CustomIntStruct struct {
	fetched int64
}

func (cis *CustomIntStruct) Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanWords(data, atEOF)
	if err == nil && token != nil {
		cis.fetched, err = strconv.ParseInt(string(token), 10, 32)
	}
	return
}

func (cis *CustomIntStruct) Fetch() int64 {
	return cis.fetched
}

func main() {
	// An artificial input source.
	//const input = "Now is the winter of our discontent,\nMade glorious summer by this sun of York.\n"
	const input = "2\na\nb\nc\nd"
	scanner := bufio.NewScanner(strings.NewReader(input))

	cis := CustomIntStruct{}
	scanner.Split(cis.Split)
	scanner.Scan()
	fmt.Println(cis.Fetch())
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	fmt.Println(scanner.Text())
	scanner.Scan()
	fmt.Println(scanner.Text())
	scanner.Scan()
	fmt.Println(scanner.Text())
	scanner.Scan()
	fmt.Println(scanner.Text())
	scanner.Scan()
	fmt.Println(scanner.Text())
}
