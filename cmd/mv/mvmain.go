package main

import (
	"fmt"
	"github.com/laher/someutils"
	"os"
)

func main() {
	err := someutils.Mv(os.Args)
	if err != nil {
		fmt.Printf("Returned error: %v\n", err)
		os.Exit(1)
	}

}
