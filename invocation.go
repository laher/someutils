package someutils

import (
	"bytes"
	"errors"
	"io"
	"os"
	"time"
)

// A set of invocation (In, Out, ErrOut, and even ErrIn (but ErrIn is usually only used by the special 'Redirector' util)
// Note that Pipables are not expected to use this type (Pipables should not need any dependency on someutils - just the implicit implementation of the Pipable interface)
type Invocation struct {
	Pipeline *Pipeline
	Pipable  Pipable
	MainPipe *Pipe
	ErrPipe  *Pipe
	/*
		MainPipe.In         io.Reader
		MainPipe.Out        io.Writer
		ErrPipe.In      io.Reader
		ErrPipe.Out     io.Writer
	*/
	SignalReceiver chan Signal
	ExitCode       *int
	Err            error
	Closed         bool
	doneChan       chan bool
}

/*
func (i *Invocation) AutoPipeErrInOut() {
	go autoPipe(i.ErrPipe.Out, i.ErrPipe.In)
}
*/
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

func (i *Invocation) WaitUpTo(timeout time.Duration) (error, *int) {
	start := time.Now()
	for {
		diff := time.Now().Sub(start)
		select {
		//todo offset timeout against time already spent in previous iterations
		case <-time.After(timeout - diff):
			return errors.New("Timeout waiting for exit codes"), nil
		case _, ok := <-i.doneChan:
			if !ok {
				return nil, i.ExitCode
			}
		}
	}
}
func (i *Invocation) Wait() *int {
	for {
		_, ok := <-i.doneChan
		if !ok {
			return i.ExitCode
		}
	}
}

/*
func (i *Invocation) PipeToPipeline(pipeline *Pipeline) (error, chan *Invocation, int) {
	if i.Pipeline != nil {
		return errors.New("This invocation already was already invoked!"), nil, 0
	}
	i.Pipeline = pipeline
	invocationChannel, count := pipeline.Invoke(i)
	return nil, invocationChannel, count
}
*/
//warning. This closes a channel. Don't call it twice - panic ensues!
func (i *Invocation) Close() error {
	if i.Closed {
		return nil
	}
	i.Closed = true
	close(i.doneChan) //this channel should only be used to check if done.
	i.MainPipe.CloseIfClosers()
	i.ErrPipe.CloseIfClosers()
	/*
	closer, ok := i.MainPipe.In.(io.Closer)
	if ok {
		err := closer.Close()
		//never mind
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not close inPipe (", err, ")")
		}
	} else {
		//		fmt.Fprintln(os.Stderr, "inPipe not a Closer")
	}
	closer, ok = i.MainPipe.Out.(io.Closer)
	if ok {
		err := closer.Close()
		//never mind
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not close outPipe (", err, ")")
		}
	} else {
		//		fmt.Fprintln(os.Stderr, "outPipe not a Closer")
	}
	closer, ok = i.ErrPipe.In.(io.Closer)
	if ok {
		err := closer.Close()
		//never mind
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not close errMainPipe.In (", err, ")")
		}
	} else {
		//		fmt.Fprintln(os.Stderr, "errMainPipe.In not a Closer")
	}
	closer, ok = i.ErrPipe.Out.(io.Closer)
	if ok {
		err := closer.Close()
		//never mind
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not close errMainPipe.Out (", err, ")")
		}
	} else {
		//		fmt.Fprintln(os.Stderr, "errMainPipe.Out not a Closer")
	}
	*/
	/*
		i.SignalReceiver <- 9
	*/
	close(i.SignalReceiver)
	return nil
}

// Convenience method returns the Stdin/Stdout/Stderr invocation associated with this process
func StdInvocation() *Invocation {
	ret := NewInvocation(os.Stdin, os.Stdout, os.Stderr)
	return ret
}

//convenience method for CLI handlers
func StdInvoke(pcu CliPipable, call []string) (error, int) {
	invocation := StdInvocation()
	err, code := pcu.ParseFlags(call, invocation.ErrPipe.Out)
	if err != nil {
		return err, code
	}
	return invocation.Pipe(pcu)
}

// Factory for a pipeline with the given invocation
func NewInvocation(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) *Invocation {
	i := new(Invocation)
	i.init()
	i.MainPipe.In = inPipe
	i.MainPipe.Out = outPipe
	i.ErrPipe.In = new(bytes.Reader)
	i.ErrPipe.Out = errPipe
	return i
}

func (i *Invocation) init() {
	i.SignalReceiver = make(chan Signal)
	i.doneChan = make(chan bool)
	i.MainPipe = new(Pipe)
	i.ErrPipe = new(Pipe)
}

// Factory taking a string for Stdin, and using byte buffers for the sdout and stderr invocation
// This returns the byte buffers to avoid the need to cast. (The assumption being that you'll want to call .Bytes() or .String() on those buffers)
func InvocationFromReader(inPipe io.Reader) (*Invocation, *bytes.Buffer, *bytes.Buffer) {
	outPipe := new(bytes.Buffer)
	errPipe := new(bytes.Buffer)
	return NewInvocation(inPipe, outPipe, errPipe), outPipe, errPipe
}
