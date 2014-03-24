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
	p := someutils.Pipeline{in, &out, &errout}
	ok, errs := p.PipeSync(Tr("H", "O"), Tr("I", "J"), Grep("O"))
	if !ok {
		fmt.Printf("Errors: %d, %+v\n", len(errs), errs)
		fmt.Printf("Errout: %+v\n", errout.String())
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}


func TestPipeline2(t *testing.T) {
	p, out, errout := someutils.NewPipelineFromString("Hi\nHo\nhI\nhO\n")
	ok, errs := p.PipeSync(Tr("H", "O"), Tr("I", "J"), Grep("O"))
	if !ok {
		t.Logf("Errors: %d, %+v\n", len(errs), errs)
		t.Logf("Errout: %+v\n", errout.String())
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}

}
