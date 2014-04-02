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
	e := p.Exec(pipes)
	err := WaitFor(e, 2, 2 * time.Second)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	// Output:
	// Hi
	// Ho
	// hI
	// hO
}
