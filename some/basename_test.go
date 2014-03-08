package some

import (
	"bytes"
	"fmt"
	"github.com/laher/someutils"
	"strings"
	"testing"
)

func TestBasename(t *testing.T) {
	var out bytes.Buffer
	var errout bytes.Buffer
	pipes := someutils.NewPipes(strings.NewReader("some/text"), &out, &errout)
	basename := NewBasename()
	err := basename.Exec(pipes)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	println(out.String())
	// TODO: 'Output' string for testing?
}
