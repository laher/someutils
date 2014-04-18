package some

import (
	"github.com/laher/someutils"
	"strings"
	"testing"
)

func TestXargsPipeline(t *testing.T) {
	pipeline := someutils.NewPipeline(Xargs(LsFact, "-l"))
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader(".\n..\n"))
	errinvocation := invocation.Wait()
	outstring := out.String()
	if errinvocation == nil {

		t.Errorf("Wait returned nil")
	}
	if errinvocation.Err != nil {
		t.Logf("errout: %+v\n", errout.String())
		t.Logf("stdout: %+v", outstring)
		t.Logf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			t.Errorf("error: %+v\n", errinvocation.Err)
		}
	}
	// TODO: 'Output' string for testing?
}
