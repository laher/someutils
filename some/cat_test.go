package some

import (
	"github.com/laher/someutils"
	"testing"
)

func TestCat(t *testing.T) {
	cat := NewCat()
	invo, outPipe, errPipe := someutils.InvocationFromString("some/text")
	err, code := cat.Invoke(invo)
	if err != nil {
		t.Logf("StdErr: %s", errPipe.String())
		t.Errorf("Error: %v - Code %d\n", err, code)
	}
	println(outPipe.String())
	// Output:
	// HI
}
