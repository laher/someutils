package main

import (
	"fmt"
	"github.com/laher/someutils/some"
	"os"
)

func main() {
	err := some.WcCli(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

}
