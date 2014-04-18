package some

import (
	"github.com/laher/someutils"
	"strings"
	"testing"
)

func TestGunzipPipeline(t *testing.T) {
	pipeline := someutils.NewPipeline(Gzip(), Gunzip(), Gzip(), Gunzip())
	input := "hello 123"
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader(input))
	invocation.Wait()
	outString := out.String()
	if invocation.Err != nil {
		t.Logf("Out: %+v\n", outString)
		t.Logf("Errout: %+v\n", errout.String())
		t.Errorf("Error: %d\n", invocation.Err)
	}
	if outString != input {
		t.Logf("Out: %+v\n", outString)
		t.Logf("Errout: %+v\n", errout.String())
		t.Errorf("Error: output doesnt match input: %s != %s\n", outString, input)
	}
}
