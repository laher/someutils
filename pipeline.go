package someutils

import (
	"bytes"
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
func (p *Pipeline) Invoke(mainInvocation *PipelineInvocation) (chan *Invocation, int) {
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
		locInvocation.init()
		if pipableIndex == 0 {
			//first iteration
			r, w = io.Pipe()
			locInvocation.MainPipe.In = mainInvocation.MainPipe.In
			locInvocation.ErrPipe.In = mainInvocation.ErrPipe.In
		} else {
			locInvocation.MainPipe.In = *previousReader
			locInvocation.ErrPipe.In = *previousErrReader
		}
		if pipableIndex == len(p.pipables)-1 {
			//last iteration
			locInvocation.MainPipe.Out = mainInvocation.MainPipe.Out
			locInvocation.ErrPipe.Out = mainInvocation.ErrPipe.Out
		} else {
			r, w = io.Pipe()
			locInvocation.MainPipe.Out = w

			rErr, wErr = io.Pipe()
			locInvocation.ErrPipe.Out = wErr
		}
		mainInvocation.Add(locInvocation)
		execAsync(pipable, locInvocation, e)
		previousReader = &r
		previousErrReader = &rErr
		pipableIndex++
	}
	return e, pipableIndex
}

func (p *Pipeline) InvokeReader(inPipe io.Reader) (*PipelineInvocation, *bytes.Buffer, *bytes.Buffer) {
	outPipe := new(bytes.Buffer)
	errPipe := new(bytes.Buffer)

	i := NewInvocation(inPipe, outPipe, errPipe)
	pi := NewPipelineInvocation(i)
	p.Invoke(pi)
	return pi, outPipe, errPipe
}

/*
// Intended as a subtype for Pipable which can redirect the error output of the previous command. This is treated as a special case because commands do not typically have access to this.
type WillRedirectErrIn interface {
	SetErrIn(errMainPipe.In io.Reader)
}
*/
// Pipe and wait for errors (up until a timeout occurs)
func (p *Pipeline) ExecAndWait(invocation *PipelineInvocation) *Invocation {
	p.Invoke(invocation)
	return invocation.Wait()
}

// Pipe and wait for errors (up until a timeout occurs)
func (p *Pipeline) ExecAndWaitUpTo(invocation *PipelineInvocation, timeout time.Duration) *Invocation {
	p.Invoke(invocation)
	return invocation.WaitUpTo(timeout)
}
