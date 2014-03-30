package some

import (
	"bytes"
	"github.com/laher/someutils"
	"strings"
	"testing"
	"time"
)

func TestExecPipeline(t *testing.T) {
	pipeline := someutils.NewPipeline(Exec("go", "help"))
	out := new(bytes.Buffer)
	errout := new(bytes.Buffer)
	ok, errs := pipeline.PipeAndWaitFor(someutils.NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), out, errout), 1 * time.Second)
	if !ok {
		t.Logf("Errors: %d, %+v\n", len(errs), errs)
		t.Logf("Errout: %+v\n", errout.String())
	}
	println(out.String())
	// TODO: 'Output' string for testing?
}
