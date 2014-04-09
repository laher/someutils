package some

import (
	"fmt"
	"github.com/laher/someutils"
	"testing"
	"time"
)

func TestExecPipeline(t *testing.T) {
	pipeline := someutils.NewPipeline(Exec("go", "help"))
	invocation, out, errout := someutils.InvocationFromString("Hi\nHo\nhI\nhO\n")
	invocationChan, count := pipeline.Invoke(invocation)

	//err, invocationchan, count := invocation.PipeToPipeline(pipeline)
/*
	if err != nil {
		fmt.Printf("error piping to pipeline: %v", err)
	}
*/
	//err, code, index := pipeline.execandwait()
	errinvocation := someutils.WaitFor(invocationChan, count, 1 * time.Second)
	outstring := out.String()
	if errinvocation == nil {
		t.Errorf("WaitFor returned nil")
	}
	if errinvocation.Err!=nil {
		fmt.Printf("errout: %+v\n", errout.String())
		fmt.Printf("stdout: %+v", outstring)
		fmt.Printf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			fmt.Printf("error: %+v\n", errinvocation.Err)
		}
	}
	fmt.Println(outstring)
	//println(out.String())
	// TODO: 'Output' string for testing?
}
