package someutils

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

type ExampleUtil struct {
}

func (ex *ExampleUtil) Exec(pipes Pipes) error {
	return LineProcessor(pipes, func(pipes Pipes, line []byte) error {
		_, err := fmt.Fprintln(pipes.Out(), string(line))
		return err
	})
}

func ExamplePipeline() {
	var errout bytes.Buffer
	//in := strings.NewReader("HiHo")
	in := strings.NewReader("Hi\nHo\nhI\nhO\n")
	pipes := NewPipes(in, os.Stdout, &errout)
	e := Pipeline(pipes, &ExampleUtil{}, &ExampleUtil{})
	errs := CollectErrors(e, 2)
	fmt.Fprintln(os.Stderr, errs)
	// Output:
	// Hi
	// Ho
	// hI
	// hO
}
