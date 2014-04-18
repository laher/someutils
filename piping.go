package someutils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func handleSignals(i *Invocation) {
	for true {
		select {
		case signal, ok := <-i.SignalReceiver:
			if !ok {
				return
			}
			quit := handleSignal(i, signal)
			if quit {
				return
			}
		}
	}
}

func handleSignal(i *Invocation, signal Signal) bool {
	switch signal.Status() {
	case 9:
		i.Close()
		return true
	case 3:
		i.Close()
		return true
	default:
		fmt.Fprintln(os.Stderr, "Unhandled signal:", signal)
		return false
	}
}

func invoke(ps PipableSimple, i *Invocation) (error, int) {
	// automatically handle the errMainPipe.In
	i.ErrPipe.Drain()
	i.AutoHandleSignals()
	//go autoPipe(i.ErrPipe.Out, i.ErrPipe.In)
	// automatically handle signals
	//go autoHandleSignals(i.signalChan, i.MainPipe.In, i.MainPipe.Out, i.ErrPipe.In, i.ErrPipe.Out)
	err, exitCode := ps.Exec(i.MainPipe.In, i.MainPipe.Out, i.ErrPipe.Out)
	//go i.Close()
	return err, exitCode
}

/*
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
*/
func autoPipe(out io.Writer, in io.Reader) {
	j, err := io.Copy(out, in)
	if err == io.EOF || err == io.ErrClosedPipe {
		//ok
		//fmt.Fprintln(os.Stderr, "expected error copying invocation", err)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "Unexpected error while copying errMainPipe.In to errMainPipe.Out", err)
	} else {
		if j > 0 {
			//fmt.Fprintln(os.Stderr, "Finished copying errMainPipe.In to errMainPipe.Out", j)
		}
		//TODO close ErrInOutinvocation here?
	}
}

// Run a Pipable asynchronously (using a goroutine)
func execAsync(pipable Pipable, invocation *Invocation, e chan *Invocation) {
	go func() {
		exitCode := execSynchronous(pipable, invocation)
		invocation.ExitCode = &exitCode
		e <- invocation
	}()
}

// run a Pipable inline
func execSynchronous(pipable Pipable, invocation *Invocation) int {
	//fmt.Fprintln(os.Stderr, "pipable starting")
	err, code := invocation.Pipe(pipable)
	//fmt.Fprintln(os.Stderr, "pipable finished")
	if err == io.EOF || err == io.ErrClosedPipe {
	} else if err != nil {
		//return &ExitError{err, code}, signalChan
		invocation.Err = err
	}
	go invocation.Close()
	return code
}

const EXIT_OK = 0

// Await completion, or first error
func Wait(e chan *Invocation, count int) *Invocation {
	var lastInvocation *Invocation

	if count < 1 {
		return NewErrorState(errors.New("No invocations to wait for!"))
	}
	for i := 0; i < count; i++ {
		select {
		case thisInvocation, ok := <-e:
			if !ok {
				fmt.Fprintln(os.Stderr, "Channel was closed!")
				break
			}
			lastInvocation = thisInvocation
			if lastInvocation.Err != nil {
				if lastInvocation.ExitCode != nil && *lastInvocation.ExitCode != EXIT_OK { //if it exited with an exitCode of OK then continue
					return lastInvocation
				}
			}
		}
	}
	return lastInvocation
}

func NewErrorState(err error) *Invocation {
	st := NewInvocation(nil, nil, nil)
	st.Err = err
	exitCode := 1
	st.ExitCode = &exitCode
	return st
}

// Await completion or error, for a duration
func WaitFor(e chan *Invocation, count int, timeout time.Duration) *Invocation {
	var lastInvocation *Invocation
	if count < 1 {
		return NewErrorState(errors.New("No invocations to wait for!"))
	}
	for i := 0; i < count; i++ {
		select {
		//todo offset timeout against time already spent in previous iterations
		case <-time.After(timeout):
			return NewErrorState(errors.New("Timeout waiting for exit codes"))
		case thisInvocation, ok := <-e:
			if !ok {
				break
			}
			lastInvocation = thisInvocation
			if lastInvocation.Err != nil {
				if lastInvocation.ExitCode != nil && *lastInvocation.ExitCode != EXIT_OK { //if it exited with an exitCode of OK then continue
					return lastInvocation
				}
			}
		}
	}
	return lastInvocation
}

// Await all errors forever
func AwaitAllErrors(e chan *Invocation, count int) (bool, []*Invocation) {
	errs := []*Invocation{}
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
func AwaitAllErrorsFor(e chan *Invocation, count int, timeout time.Duration) (bool, []*Invocation) {
	states := []*Invocation{}
	ok := true
	for i := 0; i < count; i++ {
		select {
		case <-time.After(timeout):
			states = append(states, NewErrorState(errors.New("Timeout!")))
			return false, states
		case state := <-e:
			if state.Err != nil {
				ok = false
			}
			states = append(states, state)
		}
	}
	return ok, states
}
