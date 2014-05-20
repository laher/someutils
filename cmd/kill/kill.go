package main

import (
	"fmt"
	"github.com/olekukonko/someutils"
	"os"
)

func main() {
	err := someutils.Kill(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
