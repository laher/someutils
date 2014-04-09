package some

import (
	"bytes"
	"strings"
	"testing"
)

func TestBasename(t *testing.T) {
	var outPipe bytes.Buffer
	var errPipe bytes.Buffer
	inPipe := strings.NewReader("some/text")
	basename := new(SomeBasename)
	err, code := basename.Exec(inPipe, &outPipe, &errPipe)
	if err != nil {
		t.Errorf("Error: %v - Code %d\n", err, code)
	}
	println(outPipe.String())
	// TODO: 'Output' string for testing?
}
