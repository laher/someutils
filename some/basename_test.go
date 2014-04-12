package some

import (
	"github.com/laher/someutils"
	"testing"
)

func TestBasename(t *testing.T) {
	basename := new(SomeBasename)
	invo, outPipe, errPipe := someutils.InvocationFromString("some/text")
	err, code := basename.Invoke(invo)
	if err != nil {
		t.Logf("StdErr: %s", errPipe.String())
		t.Errorf("Error: %v - Code %d\n", err, code)
	}
	println(outPipe.String())
	// TODO: 'Output' string for testing?
}
