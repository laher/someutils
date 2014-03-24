package someutils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
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
	var errout bytes.Buffer
	in := strings.NewReader("Hi\nHo\nhI\nhO\n")
	p := Pipeline{in, os.Stdout, &errout}
	e := p.Pipe(&ExampleUtil{}, &ExampleUtil{})
	errs := CollectErrors(e, 2)
	fmt.Fprintln(os.Stderr, errs)
	// Output:
	// Hi
	// Ho
	// hI
	// hO
}
