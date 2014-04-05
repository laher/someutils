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
	err, code, index := pipeline.ExecAndWaitFor(someutils.NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), out, errout), 1 * time.Second)
	if err != nil {
		t.Logf("Error: %+v, code: %d, index: %d\n", err, code, index)
		t.Logf("Errout: %+v\n", errout.String())
	}
	//println(out.String())
	// TODO: 'Output' string for testing?
}
