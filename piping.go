package someutils

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Pipes interface {
	In() io.Reader
	Out() io.Writer
	Err() io.Writer
}

type ConcretePipes struct {
	in  io.Reader
	out io.Writer
	errout io.Writer
}

func (pipes *ConcretePipes) In() io.Reader {
	return pipes.in
}

func (pipes *ConcretePipes) Out() io.Writer {
	return pipes.out
}

func (pipes *ConcretePipes) Err() io.Writer {
	return pipes.errout
}

func NewPipes(in io.Reader, out io.Writer, errout io.Writer) Pipes {
	return &ConcretePipes{in, out, errout}
}

func StdPipes() Pipes {
	return &ConcretePipes{os.Stdin, os.Stdout, os.Stderr}
}
/*
//deprecated
func Pipeline2(pipes Pipes, execable1 Execable, execable2 Execable) chan error {
	e := make(chan error)
	r, w := io.Pipe()
	pipes1 := &ConcretePipes{pipes.In(), w, pipes.Err()}
	pipes2 := &ConcretePipes{r, pipes.Out(), pipes.Err()}
	go runAsync(execable1, pipes1, w, e)
	go runAsync(execable2, pipes2, nil, e)
	return e
}
*/

func runAsync(execable Execable, pipes Pipes, closers []io.Closer, e chan error) {
	e <- execable.Exec(pipes)
	for _, closer := range closers {
		err := closer.Close()
		if err != nil {
			fmt.Fprintln(pipes.Err(), "Close error ", err)
		}
	}
}

func Pipeline(pipes Pipes, execables ...Execable) chan error {

	e := make(chan error)
	var previousReader *io.ReadCloser
	//var previousWriter *io.WriteCloser
	for i, execable := range execables {
		var w io.WriteCloser
		var r io.ReadCloser
		var in io.Reader
		var out io.Writer
		closers := []io.Closer{}
		if i == 0 {
			in = pipes.In()
		} else {
			in = *previousReader
			//closers = append (closers, *previousWriter)
			//closers = append (closers, *previousReader)
		}
		if i == len(execables)-1 {
			out = pipes.Out()
		} else {
			r, w = io.Pipe()
			out = w
			closers = append (closers, w)
		}
		go runAsync(execable, &ConcretePipes{in, out, pipes.Err()}, closers, e)
		previousReader = &r
		//previousWriter = &w
	}
	return e
}

func Collect(e chan error, count int) []error {
	errs := []error{}
	for i := 0; i < count; i++ {
		errs = append(errs, <-e)
	}
	return errs
}

type Execable interface {
	Exec(Pipes) error
}


type LineProcessorFunc func(Pipes, []byte) error

func LineProcessor(pipes Pipes, fu LineProcessorFunc) error {
	reader := bufio.NewReader(pipes.In())
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = fu(pipes, line)
		if err != nil {
			return err
		}
	}
	return nil
}

