package some

import (
	"bytes"
	"fmt"
	"github.com/laher/someutils"
	"strings"
	"testing"
)

func Test{{.NameUCF}}Pipeline(t *testing.T) {
	p, out, errout := someutils.NewPipelineFromString("Hi\nHo\nhI\nhO\n")
	ok, errs := p.PipeAndWait(1, {{.NameUCF}}("H", "O"))
	if !ok {
		t.Logf("Errors: %d, %+v\n", len(errs), errs)
		t.Logf("Errout: %+v\n", errout.String())
	}
	println(out.String())
	// TODO: 'Output' string for testing?
}
