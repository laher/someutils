package someutils

import (
	"fmt"
	"os"
	"strings"
	"bytes"
)

type ExampleUtil struct {

}

func (ex *ExampleUtil) Exec(pipes Pipes) error {
	fu := func(pipes Pipes, line []byte) error {
		_, err := fmt.Fprintln(pipes.Out(), string(line))
		return err
	}
	return LineProcessor(pipes, fu)
}


func ExamplePipeline() {
	var errout bytes.Buffer
	//in := strings.NewReader("HiHo")
	in := strings.NewReader("Hi\nHo\nhI\nhO\n")
	pipes := NewPipes(in, os.Stdout, &errout)
	e := Pipeline(pipes, &ExampleUtil{}, &ExampleUtil{})
	errs := Collect(e, 2)
	fmt.Fprintln(os.Stderr, errs)
	// Output:
	// Hi
	// Ho
	// hI
	// hO
}
