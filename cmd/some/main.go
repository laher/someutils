package main

import (
	"fmt"
	"github.com/laher/someutils"
	"os"
)
func main() {
	if len(os.Args) < 2 {
		panic("specify a command please")
	}
	err := someutils.Call(os.Args[1], os.Args[1:])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
