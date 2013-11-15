package main

import (
	"fmt"
	"github.com/laher/someutils"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "specify a command please")
		fmt.Fprintln(os.Stderr, "Commands:")
		fmt.Fprintf(os.Stderr, "  %v\n", someutils.List())
		os.Exit(1)
	}
	err := someutils.Call(os.Args[1], os.Args[1:])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
