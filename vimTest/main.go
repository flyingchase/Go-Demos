package main

import (
	"encoding/json"
	"fmt"
)

type entityData struct {
	Name string `json:"name"`
	Dep  string `json:"dep"`
	Age  int64  `json:"age"`
}

const (
	HasTop = iota + 1
	HasMiddle
	HasButton

	hasFor        = iota
	hasAnotherFor = iota
)

const hasAnother = iota + 1

const hasThird = iota + 1

func main() {
	ed := entityData{
		Name: "whoami",
		Dep:  "callMeDep",
		Age:  12132131231231,
	}
	b, err := json.Marshal(ed)
	fmt.Println(string(b))
	if err != nil {
		fmt.Printf("b: %v\n", b)
	}
	fmt.Printf("HasTop: %v\n", hasAnotherFor)
}

func quickSort(nums []int, l, r int) {
	if l > r {
		return
	}

	p := paratition(nums, l, r)
	quickSort(nums, p[1]+1, r)
	quickSort(nums, l, p[0]-1)
}

func paratition(nums []int, l, r int) []int {
	less, more := l-1, r
	for l < more {
		if nums[l] < nums[more] {

		} else if nums[l] > nums[more] {

		} else {

		}
	}
	return []int{less + 1, more}
}
