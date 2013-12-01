package main

import (
	"fmt"
	"github.com/laher/someutils"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}
	if !someutils.Exists(os.Args[1]) {
		showHelp()
		os.Exit(1)
	}
	err := someutils.Call(os.Args[1], os.Args[1:])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Fprintln(os.Stderr, "`someutils`")
	fmt.Fprintln(os.Stderr, " No command specified.")
	fmt.Fprintln(os.Stderr, " Available commands:")
	fmt.Fprintf(os.Stderr, " %v\n", someutils.List())
}
