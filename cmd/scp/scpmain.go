package main

import (
	"fmt"
	"github.com/laher/scp-go/scp"
	"os"
)

func main() {
	err := scp.Scp(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
