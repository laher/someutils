package someutils

import (
	"fmt"
	"io"
	//	"os"
	//	"time"
)

type ExampleUtil struct {
}

func (ex *ExampleUtil) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) (error, int) {
	err := LineProcessor(inPipe, outPipe, errPipe, func(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer, line []byte) error {
		_, err := fmt.Fprintln(outPipe, string(line))
		return err
	})
	if err != nil {
		return err, 1
	}
	return nil, 0
}

func ExamplePipeline() {
	p := NewPipeline(&PipableSimpleWrapper{&ExampleUtil{}}, &PipableSimpleWrapper{&ExampleUtil{}})
	mainInvocation, out, err := InvocationFromString("Hi\nHo\nhI\nhO\n") //, os.Stdout, os.Stderr)
	e, _ := p.Invoke(mainInvocation)
	<-e
	<-e
	/* errInvocation := WaitFor(e, 2, 2 * time.Second)
	if errInvocation.Err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v, Code: %d\n", errInvocation.Err, errInvocation.ExitCode)
		fmt.Fprintf(os.Stderr, "Stderr: %v\n", err.String())
	}
	*/
	fmt.Println(out.String())
	fmt.Println(err.String())
	// Output:
	// Hi
	// Ho
	// hI
	// hO
}
