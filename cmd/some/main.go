package main

import (
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/someutils/some"
	"os"
)

func main() {
	some.Init() //ensure the utils are registered.
	if len(os.Args) < 2 {
		showHelp(" No command specified.")
		os.Exit(1)
	}
	if !someutils.Exists(os.Args[1]) {
		showHelp(" Command does not exist.")
		os.Exit(1)
	}
	err := someutils.Call(os.Args[1], os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func showHelp(msg string) {
	fmt.Fprintln(os.Stderr, "`some`")
	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprintln(os.Stderr, " Available commands:")
	fmt.Fprintf(os.Stderr, " %v\n", someutils.List())
}
