package main

import (
	"fmt"
	"github.com/laher/someutils/some"
	"os"
)

func main() {
	err, code := some.RmCli(os.Args)
	if err != nil {
		if code != 0 {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(code)
		}
	}

}
