package some

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestCat(t *testing.T) {
	var out bytes.Buffer
	var errout bytes.Buffer
	inPipe, outPipe, errPipe := strings.NewReader("HI"), &out, &errout
	cat := NewCat()
	err := cat.Exec(inPipe, outPipe, errPipe)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	println(out.String())
	// Output:
	// HI
}
