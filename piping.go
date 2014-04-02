package someutils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Convenience method returns the Stdin/Stdout/Stderr pipes associated with this process
func StdPipes() (io.Reader, io.Writer, io.Writer) {
	return os.Stdin, os.Stdout, os.Stderr
}

// A set of pipes (In, Out, ErrOut, and even ErrIn (but ErrIn is usually only used by the special 'Redirector' util)
// Note that Pipables are not expected to use this type (Pipables should not need any dependency on someutils - just the implicit implementation of the Pipable interface)
type Pipeset struct {
	InPipe    io.Reader
	OutPipe   io.Writer
	ErrOutPipe   io.Writer
	ErrInPipe io.Reader
}

// Chains together the input/output of utils in a 'pipeline'
type Pipeline struct {
	pipables []Pipable
}

// Factory for a pipeline using STDIN/STDOUT/STDERR
func NewStdPipeset() *Pipeset {
	//use a bytes.Reader
	br := new(bytes.Reader)
	return &Pipeset{os.Stdin, os.Stdout, os.Stderr, br}
}

// Factory for a pipeline with the given pipes
func NewPipeset(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) *Pipeset {
	br := new(bytes.Reader)
	return &Pipeset{inPipe, outPipe, errPipe, br}
}

// Factory taking a string for Stdin, and using byte buffers for the sdout and stderr pipes
// This returns the byte buffers to avoid the need to cast. (The assumption being that you'll want to call .Bytes() or .String() on those buffers)
func NewPipesetFromString(input string) (*Pipeset, *bytes.Buffer, *bytes.Buffer) {
	var outPipe bytes.Buffer
	var errPipe bytes.Buffer
	br := new(bytes.Reader)
	return &Pipeset{strings.NewReader(input), &outPipe, &errPipe, br}, &outPipe, &errPipe
}

func NewPipeline(pipables ...Pipable) *Pipeline {
	return &Pipeline{pipables}
}


// Run a Pipable asynchronously (using a goroutine)
func execAsync(pipable Pipable, pipes *Pipeset, closers []io.Closer, e chan error) {
	rei, willRedirectErrIn := pipable.(WillRedirectErrIn)
	if willRedirectErrIn {
		rei.SetErrIn(pipes.ErrInPipe)
	} else {
		go func() {
			j, err := io.Copy(pipes.ErrOutPipe, pipes.ErrInPipe)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error! copying errInPipe to errOutPipe", err)
			}
			if j > 0 {
				fmt.Fprintln(os.Stderr, "Finished copying errInPipe to errOutPipe", j)
			}
		}()
	}
	go func() {
		e <- execSynchronous(pipable, pipes, closers)
	}()
}

// run a Pipable inline
func execSynchronous(pipable Pipable, pipes *Pipeset, closers []io.Closer) error {
	errSent := false
	if !errSent {
		err := pipable.Exec(pipes.InPipe, pipes.OutPipe, pipes.ErrOutPipe)
		if err != nil {
			return err
		}
	} else {
		//TODO show this has not run
		fmt.Fprintln(os.Stderr, "Could not run pipable.Exec")
		return errors.New("Could not run pipable.Exec")
	}
	var err error
	for _, closer := range closers {
		err = closer.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Close error ", err)
			if !errSent {
				//		return err
			}
		}
	}
	if !errSent {
		return err
	}
	return nil
}

// Run pipables in a sequence, weaving together their inputs and outputs appropriately
func (p *Pipeline) Exec(pipes *Pipeset) chan error {
	e := make(chan error)
	var previousReader *io.ReadCloser
	var previousErrReader *io.ReadCloser
	//fmt.Printf("%+v\n", pipables)
	for i, pipable := range p.pipables {
		//println(pipable)
		var w io.WriteCloser
		var r io.ReadCloser
		var wErr io.WriteCloser
		var rErr io.ReadCloser
		locpipes := new(Pipeset)
		closers := []io.Closer{}
		if i == 0 {
			locpipes.InPipe = pipes.InPipe
			locpipes.ErrInPipe = pipes.ErrInPipe
		} else {
			locpipes.InPipe = *previousReader
			locpipes.ErrInPipe = *previousErrReader
		}
		if i == len(p.pipables)-1 {
			locpipes.OutPipe = pipes.OutPipe

			outCloser, isCloser := locpipes.OutPipe.(io.Closer)
			if isCloser {
				closers = append(closers, outCloser)
			}
			locpipes.ErrOutPipe = pipes.ErrOutPipe
			errOutCloser, isCloser := locpipes.ErrOutPipe.(io.Closer)
			if isCloser {
				closers = append(closers, errOutCloser)
			}
		} else {
			r, w = io.Pipe()
			locpipes.OutPipe = w
			closers = append(closers, w)

			rErr, wErr = io.Pipe()
			locpipes.ErrOutPipe = wErr
			closers = append(closers, wErr)
		}
		execAsync(pipable, locpipes, closers, e)
		previousReader = &r
		previousErrReader = &rErr
	}
	return e
}

// Intended as a subtype for Pipable which can redirect the error output of the previous command. This is treated as a special case because commands do not typically have access to this.
type WillRedirectErrIn interface {
	SetErrIn(errInPipe io.Reader)
}

// Pipe and wait for errors (up until a timeout occurs)
func (p *Pipeline) ExecAndWait(pipes *Pipeset) error {
	e := p.Exec(pipes)
	return Wait(e, len(p.pipables))
}


// Pipe and wait for errors (up until a timeout occurs)
func (p *Pipeline) ExecAndWaitFor(pipes *Pipeset, timeout time.Duration) error {
	e := p.Exec(pipes)
	return WaitFor(e, len(p.pipables), timeout)
}


// Await first error
func Await(e chan error, count int) error {
	for i := 0; i < count; i++ {
		select {
		case err := <-e:
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Await completion, or first error
func Wait(e chan error, count int) error {
	for i := 0; i < count; i++ {
		select {
		case err := <-e:
			if err != nil {
				return err
			}
		}
	}
	return nil
}


// Await completion or error, for a duration
func WaitFor(e chan error, count int, timeout time.Duration) error {
	for i := 0; i < count; i++ {
		select {
		//todo offset timeout against time already spent in previous iterations
		case <-time.After(timeout):
			return errors.New("Timeout!")
		case err := <-e:
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Await all errors forever
func AwaitAllErrors(e chan error, count int) (bool, []error) {
	errs := []error{}
	ok := true
	for i := 0; i < count; i++ {
		select {
		case err := <-e:
			if err != nil {
				ok = false
			}
			errs = append(errs, err)
		}
	}
	return ok, errs
}

// Await Errors for a duration
func AwaitAllErrorsFor(e chan error, count int, timeout time.Duration) (bool, []error) {
	errs := []error{}
	ok := true
	for i := 0; i < count; i++ {
		select {
		case <-time.After(timeout):
			errs = append(errs, errors.New("Timeout!"))
			return false, errs
		case err := <-e:
			if err != nil {
				ok = false
			}
			errs = append(errs, err)
		}
	}
	return ok, errs
}

// A function which processses a line from a reader
type LineProcessorFunc func(io.Reader, io.Writer, io.Writer, []byte) error

// Process line-by-line
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
