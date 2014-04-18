package main

import (
	"fmt"
	"github.com/laher/someutils/some"
	"os"
)

func main() {
	err := some.WgetCli(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

}
