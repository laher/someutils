package some

import (
	"github.com/laher/someutils"
	"strings"
	"testing"
)

func Test{{.NameUCF}}Pipeline(t *testing.T) {
	pipeline := someutils.NewPipeline({{.NameUCF}}("H", "O"))
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader("Hi\nHo\nhI\nhO\n"))
	invocation.Wait()
	t.Logf("Out: %+v\n", out.String())
	t.Logf("Errout: %+v\n", errout.String())
	if invocation.Err != nil {
		t.Errorf("Error: %d\n", err)
	}
	// TODO: 'Output' string for testing?
}
