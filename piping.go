package someutils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type Execable interface {
	Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error
}

func StdPipes() (io.Reader, io.Writer, io.Writer) {
	return os.Stdin, os.Stdout, os.Stderr
}

type Pipeline struct {
	InPipe  io.Reader
	OutPipe io.Writer
	ErrPipe io.Writer
}

func NewStdPipeline() *Pipeline {
	return &Pipeline{os.Stdin, os.Stdout, os.Stderr}
}

func NewPipelineFromString(input string) (*Pipeline, *bytes.Buffer, *bytes.Buffer) {
	var outPipe bytes.Buffer
	var errPipe bytes.Buffer
	return &Pipeline{strings.NewReader(input), &outPipe, &errPipe}, &outPipe, &errPipe
}

func runAsync(execable Execable, inPipe io.Reader, outPipe io.Writer, errPipe io.Writer, closers []io.Closer, e chan error) {
	e <- execable.Exec(inPipe, outPipe, errPipe)
	for _, closer := range closers {
		err := closer.Close()
		if err != nil {
			fmt.Fprintln(errPipe, "Close error ", err)
		}
	}
}

func (p *Pipeline) Pipe(execables ...Execable) chan error {
	e := make(chan error)
	var previousReader *io.ReadCloser
	for i, execable := range execables {
		var w io.WriteCloser
		var r io.ReadCloser
		var locInPipe io.Reader
		var locOutPipe io.Writer
		closers := []io.Closer{}
		if i == 0 {
			locInPipe = p.InPipe
		} else {
			locInPipe = *previousReader
		}
		if i == len(execables)-1 {
			locOutPipe = p.OutPipe
		} else {
			r, w = io.Pipe()
			locOutPipe = w
			closers = append(closers, w)
		}
		go runAsync(execable, locInPipe, locOutPipe, p.ErrPipe, closers, e)
		previousReader = &r
	}
	return e
}

func (p *Pipeline) PipeSync(execables ...Execable) []error {
	e := p.Pipe(execables...)
	return CollectErrors(e, len(execables))
}

func CollectErrors(e chan error, count int) []error {
	errs := []error{}
	for i := 0; i < count; i++ {
		errs = append(errs, <-e)
	}
	return errs
}

type LineProcessorFunc func(io.Reader, io.Writer, io.Writer, []byte) error

func LineProcessor(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer, fu LineProcessorFunc) error {
	reader := bufio.NewReader(inPipe)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = fu(inPipe, outPipe, errPipe, line)
		if err != nil {
			return err
		}
	}
}
