package main

import (
	"fmt"
	wgetgo "github.com/laher/wget-go"
	"os"
)

func main() {
	err := wgetgo.Wget(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

}
