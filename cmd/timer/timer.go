package main

import (
	"fmt"
	"github.com/olekukonko/someutils"
	"os"
)

func main() {
	err := someutils.Timer(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
