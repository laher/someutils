//+build freebsd openbsd netbsd

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintf(os.Stderr, "Error: %v scp uses gopass which is not suported in *bsd except via c calls\n")
}
