package some

import (
	"bytes"
	"fmt"
	"github.com/laher/someutils"
	"strings"
	"testing"
)

func TestCat(t *testing.T) {
	var out bytes.Buffer
	var errout bytes.Buffer
	pipes := someutils.NewPipes(strings.NewReader("HI"), &out, &errout)
	cat := NewCat()
	err := cat.Exec(pipes)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	println(out.String())
	// Output:
	// HI
}
