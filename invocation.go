package someutils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// A set of invocation (In, Out, ErrOut, and even ErrIn (but ErrIn is usually only used by the special 'Redirector' util)
// Note that Pipables are not expected to use this type (Pipables should not need any dependency on someutils - just the implicit implementation of the Pipable interface)
type Invocation struct {
	Pipeline       *Pipeline
	Pipable        Pipable
	InPipe         io.Reader
	OutPipe        io.Writer
	ErrInPipe      io.Reader
	ErrOutPipe     io.Writer
	SignalReceiver chan int
	ExitCode       *int
	Err            error
	Closed         bool
}

func (i *Invocation) AutoPipeErrInOut() {
	go autoPipe(i.ErrOutPipe, i.ErrInPipe)
}

//TODO!!
func (i *Invocation) AutoHandleSignals() {
	//	go handleSignals(i)
}

func (i *Invocation) Pipe(pipable Pipable) (error, int) {
	if i.Pipable != nil {
		return errors.New("This invocation already was already invoked!"), 1
	}
	i.Pipable = pipable
	return pipable.Invoke(i)
}

func (i *Invocation) PipeToPipeline(pipeline *Pipeline) (error, chan *Invocation, int) {
	if i.Pipeline != nil {
		return errors.New("This invocation already was already invoked!"), nil, 0
	}
	i.Pipeline = pipeline
	invocationChannel, count := pipeline.Invoke(i)
	return nil, invocationChannel, count
}

//warning. This closes a channel. Don't call it twice - panic ensues!
func (i *Invocation) Close() error {
	if i.Closed {
		return nil
	}
	i.Closed = true
	closer, ok := i.InPipe.(io.Closer)
	if ok {
		err := closer.Close()
		//never mind
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not close inPipe (", err, ")")
		}
	} else {
		//		fmt.Fprintln(os.Stderr, "inPipe not a Closer")
	}
	closer, ok = i.OutPipe.(io.Closer)
	if ok {
		err := closer.Close()
		//never mind
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not close outPipe (", err, ")")
		}
	} else {
		//		fmt.Fprintln(os.Stderr, "outPipe not a Closer")
	}
	closer, ok = i.ErrInPipe.(io.Closer)
	if ok {
		err := closer.Close()
		//never mind
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not close errInPipe (", err, ")")
		}
	} else {
		//		fmt.Fprintln(os.Stderr, "errInPipe not a Closer")
	}
	closer, ok = i.ErrOutPipe.(io.Closer)
	if ok {
		err := closer.Close()
		//never mind
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not close errOutPipe (", err, ")")
		}
	} else {
		//		fmt.Fprintln(os.Stderr, "errOutPipe not a Closer")
	}
	/*
		i.SignalReceiver <- 9
		close(i.SignalReceiver)
	*/
	return nil
}

// Convenience method returns the Stdin/Stdout/Stderr invocation associated with this process
func StdInvocation() *Invocation {
	ps := new(Invocation)
	ps.InPipe = os.Stdin
	ps.OutPipe = os.Stdout
	ps.ErrOutPipe = os.Stderr
	ps.ErrInPipe = strings.NewReader("")
	ps.SignalReceiver = make(chan int)
	return ps
}

//convenience method for CLI handlers
func StdInvoke(pcu CliPipable, call []string) (error, int) {
	invocation := StdInvocation()
	err, code := pcu.ParseFlags(call, invocation.ErrOutPipe)
	if err != nil {
		return err, code
	}
	return invocation.Pipe(pcu)
}

// Factory for a pipeline with the given invocation
func NewInvocation(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) *Invocation {
	i := new(Invocation)
	i.InPipe = inPipe
	i.OutPipe = outPipe
	i.ErrInPipe = new(bytes.Reader)
	i.ErrOutPipe = errPipe
	i.SignalReceiver = make(chan int)
	return i
}

// Factory taking a string for Stdin, and using byte buffers for the sdout and stderr invocation
// This returns the byte buffers to avoid the need to cast. (The assumption being that you'll want to call .Bytes() or .String() on those buffers)
func InvocationFromString(input string) (*Invocation, *bytes.Buffer, *bytes.Buffer) {
	var outPipe bytes.Buffer
	var errPipe bytes.Buffer
	return NewInvocation(strings.NewReader(input), &outPipe, &errPipe), &outPipe, &errPipe
}
