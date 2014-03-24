package some

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestBasename(t *testing.T) {
	var outPipe bytes.Buffer
	var errPipe bytes.Buffer
	inPipe := strings.NewReader("some/text")
	basename := NewBasename()
	err := basename.Exec(inPipe, &outPipe, &errPipe)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	println(outPipe.String())
	// TODO: 'Output' string for testing?
}
