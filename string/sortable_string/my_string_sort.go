package sortable_string

import (
	"fmt"
	"sort"
)

type StringArray []string

func (input StringArray) Len() int {
	return len(input)
}

func (input StringArray) Less(i, j int) bool {
	return input[i] < input[j]
}

func (input StringArray) Swap(i, j int) {
	input[i], input[j] = input[j], input[i]
}

/*func main() {
	ma := StringArray{"8","7","6"}
	fmt.Println(ma)
	sort.Sort(ma)
	fmt.Println(ma)
}*/