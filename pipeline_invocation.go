package someutils

import (
	"time"
	//"fmt"
)

type PipelineInvocation struct {
	*Invocation
	Invocations []*Invocation
}

func (pi *PipelineInvocation) Add(invocation *Invocation) {
	if pi.Invocations == nil {
		pi.Invocations = []*Invocation{invocation}
	} else {
		pi.Invocations = append(pi.Invocations, invocation)
	}
}

func (pi *PipelineInvocation) SignalAll(signal Signal) {
	for _, i := range pi.Invocations {
		i.SignalReceiver <- signal
	}
}

// Wait until an exit status has occurred
func (pi *PipelineInvocation) Wait() *Invocation {
	var i *Invocation
	//fmt.Printf("invocations: %v\n", pi.Invocations)
	for _, i = range pi.Invocations {
		statusCode := i.Wait()
		if statusCode == nil || *statusCode != 0 {
			return i
		}
	}
	return i //last one (or nil if no invocations)
}

// Wait until an exit status has occurred
func (pi *PipelineInvocation) WaitUpTo(timeout time.Duration) *Invocation {
	var i *Invocation
	start := time.Now()
	for _, i = range pi.Invocations {
		diff := time.Now().Sub(start)
		err, statusCode := i.WaitUpTo(timeout - diff)
		if err != nil {
			return NewErrorState(err)
		}
		if statusCode == nil || *statusCode != 0 {
			return i
		}
	}
	return i //last one (or nil if no invocations)
}

func NewPipelineInvocation(invocation *Invocation) *PipelineInvocation {
	pi := new(PipelineInvocation)
	pi.Invocation = invocation
	return pi
}

/*
// This returns the byte buffers to avoid the need to cast. (The assumption being that you'll want to call .Bytes() or .String() on those buffers)
func PipelineInvocationFromReader(inPipe io.Reader) (*PipelineInvocation, *bytes.Buffer, *bytes.Buffer) {
	i, o, e :=  InvocationFromReader(inPipe)
	pi := new(PipelineInvocation)
	pi.Invocation = i
	return pi, o, e
}
*/
