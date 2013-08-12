package main

import (
	"github.com/laher/someutils"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("specify a command please")
	}
	switch os.Args[1] {
		case "cp" :
			someutils.Cp(os.Args[1:])
		case "ls" :
			someutils.Ls(os.Args[1:])
		case "mv" :
			someutils.Mv(os.Args[1:])
		case "rm" :
			someutils.Rm(os.Args[1:])
	}
}
