package some

import (
	"bytes"
	"fmt"
	"github.com/laher/someutils"
	"strings"
	"testing"
)

func Test{{.NameUCF}}(t *testing.T) {
	var out bytes.Buffer
	var errout bytes.Buffer
	pipes := someutils.NewPipes(strings.NewReader("some/text"), &out, &errout)
	{{.Name}} := New{{.NameUCF}}()
	err := {{.Name}}.Exec(pipes)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	println(out.String())
	// TODO: 'Output' string for testing?
}
