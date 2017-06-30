package template

import (
	"strconv"
	"strings"
	"bufio"
	"io"
)

////////////////////////////////////////////////////////////////////////////////

// INPUT TEMPLATE START

type MyInput struct {
	rdr         io.Reader
	lineChan    chan string
	initialized bool
}

func NewMyInput(rdr *strings.Reader) *MyInput {
	return &MyInput{rdr: rdr}
}

func (mi *MyInput) start(done chan struct{}) {
	r := bufio.NewReader(mi.rdr)
	defer func() { close(mi.lineChan) }()
	for {
		line, err := r.ReadString('\n')
		if !mi.initialized {
			mi.initialized = true
			done <- struct{}{}
		}
		mi.lineChan <- strings.TrimSpace(line)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}
}

func (mi *MyInput) ReadWords() []string {
	line := mi.readline()
	return strings.Split(line, " ")
}

func (mi *MyInput) readline() string {
	// if this is the first call, initialize
	if !mi.initialized {
		mi.lineChan = make(chan string)
		done := make(chan struct{})
		go mi.start(done)
		<-done
	}

	res, ok := <-mi.lineChan
	if !ok {
		panic("trying to read from a closed channel")
	}
	return res
}

func (mi *MyInput) ReadInt() int {
	line := mi.readline()
	i, err := strconv.Atoi(line)
	if err != nil {
		panic(err)
	}
	return i
}

func (mi *MyInput) ReadInt64() int64 {
	line := mi.readline()
	i, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (mi *MyInput) ReadInts() []int {
	line := mi.readline()
	parts := strings.Split(line, " ")
	res := []int{}
	for _, s := range parts {
		tmp, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		res = append(res, tmp)
	}
	return res
}

func (mi *MyInput) ReadInt64s() []int64 {
	line := mi.readline()
	parts := strings.Split(line, " ")
	res := []int64{}
	for _, s := range parts {
		tmp, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(err)
		}
		res = append(res, tmp)
	}
	return res
}

// INPUT TEMPLATE END

////////////////////////////////////////////////////////////////////////////////
