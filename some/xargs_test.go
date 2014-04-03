package some

import (
	"bytes"
	"github.com/laher/someutils"
	"strings"
	"testing"
)

func TestXargsPipeline(t *testing.T) {
	pipeline := someutils.NewPipeline(Xargs(LsFactory, "-l"))
	var out bytes.Buffer
	var errout bytes.Buffer
	in := strings.NewReader(".\n..\n")
	err := pipeline.ExecAndWait(someutils.NewPipeset(in, &out, &errout))
	t.Logf("Out (length): %d\n", len(out.String()))
	t.Logf("Errout: %+v\n", errout.String())
	if err != nil {
		t.Errorf("Error: %d\n", err)
	}
	// TODO: 'Output' string for testing?
}
