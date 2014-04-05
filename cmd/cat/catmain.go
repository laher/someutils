package main

import (
	"fmt"
	"github.com/laher/someutils/some"
	"os"
)

func main() {
	err, code := some.CatCli(os.Args)
	if err != nil {
		if code != 0 {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(code)
		}
	}

}
