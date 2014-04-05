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

func (p *Pipeset) Close() error {
	closer, ok := p.InPipe.(io.Closer)
	if ok {
		err := closer.Close()
		//never mind
		if err!=nil {
			fmt.Fprintln(os.Stderr, "Could not close inPipe (", err, ")")
		}
	} else {
//		fmt.Fprintln(os.Stderr, "inPipe not a Closer")
	}
	closer, ok = p.OutPipe.(io.Closer)
	if ok {
		err := closer.Close()
		//never mind
		if err!= nil {
			fmt.Fprintln(os.Stderr, "Could not close outPipe (", err, ")")
		}
	} else {
//		fmt.Fprintln(os.Stderr, "outPipe not a Closer")
	}
	closer, ok = p.ErrInPipe.(io.Closer)
	if ok {
		err := closer.Close()
		//never mind
		if err!= nil {
			fmt.Fprintln(os.Stderr, "Could not close errInPipe (", err, ")")
		}
	} else {
//		fmt.Fprintln(os.Stderr, "errInPipe not a Closer")
	}
	closer, ok = p.ErrOutPipe.(io.Closer)
	if ok {
		err := closer.Close()
		//never mind
		if err!= nil {
			fmt.Fprintln(os.Stderr, "Could not close errOutPipe (", err, ")")
		}
	} else {
//		fmt.Fprintln(os.Stderr, "errOutPipe not a Closer")
	}
	return nil
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

type CodeError interface {
	Err() error
	Code() int
}

type ExitError struct {
	err error
	ExitCode int
}

func (exitError *ExitError) Err() error {
	return exitError.err
}
func (exitError *ExitError) Code() int {
	return exitError.ExitCode
}

// Run a Pipable asynchronously (using a goroutine)
func execAsync(pipable Pipable, pipes *Pipeset, e chan CodeError) {
	rei, willRedirectErrIn := pipable.(WillRedirectErrIn)
	if willRedirectErrIn {
		rei.SetErrIn(pipes.ErrInPipe)
	} else {
		go func() {
			j, err := io.Copy(pipes.ErrOutPipe, pipes.ErrInPipe)
			if err == io.EOF || err == io.ErrClosedPipe {
				//ok
				//fmt.Fprintln(os.Stderr, "expected error copying pipes", err)
			} else if err != nil {
				fmt.Fprintln(os.Stderr, "Unexpected error while copying errInPipe to errOutPipe", err)
			} else {
				if j > 0 {
					//fmt.Fprintln(os.Stderr, "Finished copying errInPipe to errOutPipe", j)
				}
				//TODO close ErrInOutPipes here?
			}
		}()
	}
	go func() {
		e <- execSynchronous(pipable, pipes)
	}()
}

// run a Pipable inline
func execSynchronous(pipable Pipable, pipes *Pipeset) CodeError {
	//fmt.Fprintln(os.Stderr, "pipable starting")
	err, code := pipable.Exec(pipes.InPipe, pipes.OutPipe, pipes.ErrOutPipe)
	//fmt.Fprintln(os.Stderr, "pipable finished")
	pipes.Close()
	if err == io.EOF || err == io.ErrClosedPipe {
		return nil
	} else if err != nil {
		return &ExitError{err, code}
	}
	return nil
}

// Run pipables in a sequence, weaving together their inputs and outputs appropriately
func (p *Pipeline) Exec(pipes *Pipeset) chan CodeError {
	e := make(chan CodeError)
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
		if i == 0 {
			//first iteration
			r, w = io.Pipe()
			locpipes.InPipe = pipes.InPipe
			locpipes.ErrInPipe = pipes.ErrInPipe
		} else {
			locpipes.InPipe = *previousReader
			locpipes.ErrInPipe = *previousErrReader
		}
		if i == len(p.pipables)-1 {
			//last iteration
			locpipes.OutPipe = pipes.OutPipe
			locpipes.ErrOutPipe = pipes.ErrOutPipe
		} else {
			r, w = io.Pipe()
			locpipes.OutPipe = w

			rErr, wErr = io.Pipe()
			locpipes.ErrOutPipe = wErr
		}
		execAsync(pipable, locpipes, e)
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
func (p *Pipeline) ExecAndWait(pipes *Pipeset) (error, int, int) {
	e := p.Exec(pipes)
	return Wait(e, len(p.pipables))
}


// Pipe and wait for errors (up until a timeout occurs)
func (p *Pipeline) ExecAndWaitFor(pipes *Pipeset, timeout time.Duration) (error, int, int) {
	e := p.Exec(pipes)
	return WaitFor(e, len(p.pipables), timeout)
}


// Await completion, or first error
func Wait(e chan CodeError, count int) (error, int, int) {
	i := 0
	for ; i < count; i++ {
		select {
		case codeError := <-e:
			if codeError != nil {
				return codeError.Err(), codeError.Code(), i
			}
		}
	}
	return nil, 0, i
}


// Await completion or error, for a duration
func WaitFor(e chan CodeError, count int, timeout time.Duration) (error, int, int) {
	i := 0
	for ; i < count; i++ {
		select {
		//todo offset timeout against time already spent in previous iterations
		case <-time.After(timeout):
			return errors.New("Timeout!"), 1, i
		case codeError := <-e:
			if codeError != nil {
				return codeError.Err(), codeError.Code(), i
			}
		}
	}
	return nil, 0, i
}

// Await all errors forever
func AwaitAllErrors(e chan CodeError, count int) (bool, []CodeError) {
	errs := []CodeError{}
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
func AwaitAllErrorsFor(e chan CodeError, count int, timeout time.Duration) (bool, []CodeError) {
	errs := []CodeError{}
	ok := true
	for i := 0; i < count; i++ {
		select {
		case <-time.After(timeout):
			errs = append(errs, &ExitError{errors.New("Timeout!"), 1})
			return false, errs
		case err := <-e:
			if err.Err() != nil {
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
		//println("get line")
		line, _, err := reader.ReadLine()
		//println("got line", line)
		if err == io.EOF || err == io.ErrClosedPipe {
			return nil
		}
		if err != nil {
			return err
		}
		err = fu(inPipe, outPipe, errPipe, line)
		//println("ran fu", line)
		if err != nil {
			return err
		}
	}
}
