package someutils

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type ExampleUtil struct {
}

func (ex *ExampleUtil) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	return LineProcessor(inPipe, outPipe, errPipe, func(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer, line []byte) error {
		_, err := fmt.Fprintln(outPipe, string(line))
		return err
	})
}

func ExamplePipeline() {
	p := NewPipeline(&ExampleUtil{}, &ExampleUtil{})
	pipes := NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), os.Stdout, os.Stderr)
	e := p.Pipe(pipes)
	ok, errs := AwaitErrorsFor(e, 2, 2 * time.Second)
	if !ok {
		fmt.Fprintln(os.Stderr, errs)
	}
	// Output:
	// Hi
	// Ho
	// hI
	// hO
}
