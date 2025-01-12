package main

import (
	"fmt"
	"os"
)

func main() {
	// TODO: actually it should check directory too, ~/lab/.lab
	Setup()

	firstArg := os.Args[0]
	if len(firstArg) < 1 {
		fmt.Println("no arguments")
	}
	fmt.Printf("First ar is %v", firstArg)
}
