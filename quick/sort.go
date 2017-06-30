package main

import (
	"fmt"
	"sort"
)

type IntArray []int

func (input IntArray) Len() int {
	return len(input)
}

func (input IntArray) Less(i, j int) bool {
	return input[i] < input[j]
}

func (input IntArray) Swap(i, j int) {
	input[i], input[j] = input[j], input[i]
}

func main() {
	ma := IntArray{8,7,6}
	fmt.Println(ma)
	sort.Sort(ma)
	fmt.Println(ma)
}