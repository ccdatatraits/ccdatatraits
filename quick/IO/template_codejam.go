package main

import (
	"strings"
	"fmt"
	"github.com/ccdatatraits/template"
)

func main() {
	const input = "2\na b\nc d"

	mi := template.NewMyInput(/*os.Stdin*/strings.NewReader(input))
	t := mi.ReadInt()
	for caseNo := 1; caseNo <= t; caseNo++ {
		words := mi.ReadWords()
		fmt.Println(caseNo, ":", strings.Join(words, "-"))
	}
}
