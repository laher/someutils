package some

import (
	"github.com/laher/someutils"
	"testing"
)

func TestXargsPipeline(t *testing.T) {
	pipeline := someutils.NewPipeline(Xargs(LsFact, "-l"))
	invocation, out, errout := someutils.InvocationFromString(".\n..\n")
	err, invocationchan, count := invocation.PipeToPipeline(pipeline)
	if err != nil {
		t.Errorf("error piping to pipeline: %v", err)
	}
	errinvocation := someutils.Wait(invocationchan, count)
	outstring := out.String()
	if errinvocation == nil {

			t.Errorf("Wait returned nil. Expecting %d invocations", count)
	}
	if errinvocation.Err!=nil {
		t.Logf("errout: %+v\n", errout.String())
		t.Logf("stdout: %+v", outstring)
		t.Logf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			t.Errorf("error: %+v\n", errinvocation.Err)
		}
	}
	// TODO: 'Output' string for testing?
}
