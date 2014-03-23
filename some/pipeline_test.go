package some

import (
	"bytes"
	"fmt"
	"github.com/laher/someutils"
	"strings"
	"testing"
)

func TestPipeline1(t *testing.T) {
	var out bytes.Buffer
	var errout bytes.Buffer
	in := strings.NewReader("Hi\nHo\nhI\nhO\n")
	pipes := someutils.NewPipes(in, &out, &errout)
	errs := someutils.PipelineSync(pipes, Tr("H", "O"), Tr("I", "J"), Grep("O"))
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
	fmt.Printf("Errors: %d, %+v\n", len(errs), errs)
	fmt.Printf("Errout: %+v\n", errout.String())
}
