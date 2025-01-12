package main

import (
	"fmt"
	"os"
)

func main() {
	// TODO: actually it should check directory too, ~/lab/.lab
	ensureConfigFile()
	firstArg := os.Args[1]
	// if len(firstArg) < 1 {
	// 	fmt.Println("no arguments")
	// }
	fmt.Printf("First ar is %v", firstArg)
}
