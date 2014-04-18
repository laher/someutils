package main

import (
	"fmt"
	"github.com/laher/wget-go/wget"
	"os"
)

func main() {
	err, code := wget.WgetCli(os.Args)
	if err != nil {
		if code != 0 {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(code)
		}
	}

}
