package some

import (
	"fmt"
	"github.com/laher/someutils"
	"strings"
	"testing"
	"time"
)

func TestExecPipeline(t *testing.T) {
	pipeline := someutils.NewPipeline(Exec("go", "help"))
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader("Hi\nHo\nhI\nhO\n"))
	errinvocation := invocation.WaitUpTo(1 * time.Second)
	outstring := out.String()
	if errinvocation == nil {
		t.Errorf("WaitFor returned nil")
	}
	if errinvocation.Err != nil {
		fmt.Printf("errout: %+v\n", errout.String())
		fmt.Printf("stdout: %+v", outstring)
		fmt.Printf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			fmt.Printf("error: %+v\n", errinvocation.Err)
		}
	}
	// TODO: 'Output' string for testing?
	fmt.Println(outstring[0:10])
	// Output:
	// Go is a to
}
