//+build freebsd openbsd netbsd

package main


import (
	"fmt"
)

func main() {
	fmt.Printf("Error: %v scp uses gopass which is not suported in *bsd except via c calls\n")
}
