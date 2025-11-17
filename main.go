package main

import (
	"fmt"

	"github.com/mykeelium/visual-playground/collatz"
)

func main() {
	tree := collatz.BuildTree(100000)

	fmt.Println("tree:")
	collatz.PrintOrganicTree(&tree)
}
