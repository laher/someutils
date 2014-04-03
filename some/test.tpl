package some

import (
	"bytes"
	"fmt"
	"github.com/laher/someutils"
	"strings"
	"testing"
)

func Test{{.NameUCF}}Pipeline(t *testing.T) {
	pipeline := someutils.NewPipeline({{.NameUCF}}("H", "O"))
	out := new(bytes.Buffer)
	errout := new(bytes.Buffer)
	err := pipeline.ExecAndWaitFor(someutils.NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), out, errout), 1 * time.Second)
	t.Logf("Out: %+v\n", out.String())
	t.Logf("Errout: %+v\n", errout.String())
	if err != nil {
		t.Errorf("Error: %d\n", err)
	}
	// TODO: 'Output' string for testing?
}
