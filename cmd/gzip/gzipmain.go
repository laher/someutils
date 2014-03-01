package main

import (
	"fmt"
	"github.com/laher/someutils"
	"os"
)

func main() {
	err := someutils.Gzip(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
