package someutils

import (
	"io"
	"time"
)



// Chains together the input/output of utils in a 'pipeline'
type Pipeline struct {
	pipables []Pipable
}



func NewPipeline(pipables ...Pipable) *Pipeline {
	return &Pipeline{pipables}
}


// Run pipables in a sequence, weaving together their inputs and outputs appropriately
func (p *Pipeline) Invoke(mainInvocation *Invocation) (chan *Invocation, int) {
	e := make(chan *Invocation)
	var previousReader *io.ReadCloser
	var previousErrReader *io.ReadCloser
	//fmt.Printf("%+v\n", pipables)
	pipableIndex := 0
	for _, pipable := range p.pipables {
		//println(pipable)
		var w io.WriteCloser
		var r io.ReadCloser
		var wErr io.WriteCloser
		var rErr io.ReadCloser
		locInvocation := new(Invocation)
		if pipableIndex == 0 {
			//first iteration
			r, w = io.Pipe()
			locInvocation.InPipe = mainInvocation.InPipe
			locInvocation.ErrInPipe = mainInvocation.ErrInPipe
		} else {
			locInvocation.InPipe = *previousReader
			locInvocation.ErrInPipe = *previousErrReader
		}
		if pipableIndex == len(p.pipables)-1 {
			//last iteration
			locInvocation.OutPipe = mainInvocation.OutPipe
			locInvocation.ErrOutPipe = mainInvocation.ErrOutPipe
		} else {
			r, w = io.Pipe()
			locInvocation.OutPipe = w

			rErr, wErr = io.Pipe()
			locInvocation.ErrOutPipe = wErr
		}
		execAsync(pipable, locInvocation, e)
		previousReader = &r
		previousErrReader = &rErr
		pipableIndex++
	}
	return e, pipableIndex
}
/*
// Intended as a subtype for Pipable which can redirect the error output of the previous command. This is treated as a special case because commands do not typically have access to this.
type WillRedirectErrIn interface {
	SetErrIn(errInPipe io.Reader)
}
*/
// Pipe and wait for errors (up until a timeout occurs)
func (p *Pipeline) ExecAndWait(invocation *Invocation) *Invocation {
	e, count := p.Invoke(invocation)
	return Wait(e, count)
}


// Pipe and wait for errors (up until a timeout occurs)
func (p *Pipeline) ExecAndWaitFor(invocation *Invocation, timeout time.Duration) *Invocation {
	e, count := p.Invoke(invocation)
	return WaitFor(e, count, timeout)
}
