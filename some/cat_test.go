package some

import (
	"bytes"
	"strings"
	"testing"
)

func TestCat(t *testing.T) {
	var out bytes.Buffer
	var errout bytes.Buffer
	inPipe, outPipe, errPipe := strings.NewReader("HI"), &out, &errout
	cat := NewCat()
	err, code := cat.Exec(inPipe, outPipe, errPipe)
	if err != nil {
		t.Errorf("Error: %v - Code %d\n", err, code)
	}
	println(out.String())
	// Output:
	// HI
}
