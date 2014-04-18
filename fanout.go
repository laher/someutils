package someutils

import (
	"bufio"
	"strings"
	"sync"
)


// SomeFanout represents and performs a `fanout` invocation
type SomeFanout struct {
	// TODO: maxRoutines int
	strategy FanoutStrategy
	pipeline *Pipeline
}

type FanoutStrategy func(*SomeFanout, *Invocation) (chan *PipelineInvocation, chan error)

func FanoutByLineStrategy (fanout *SomeFanout, invocation *Invocation) (chan *PipelineInvocation, chan error) {
	piChan := make(chan *PipelineInvocation)
	errChan := make(chan error)
	go func() {
		reader := bufio.NewReader(invocation.MainPipe.In)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				errChan <- err
				close(piChan)
				return
			}
			inv := NewInvocation(strings.NewReader(string(line)), invocation.MainPipe.Out, invocation.ErrPipe.Out)
			pi := NewPipelineInvocation(inv)
			fanout.pipeline.Invoke(pi)
			piChan <- pi
		}
	}()
	return piChan, errChan
}

// Name() returns the name of the util
func (fanout *SomeFanout) Name() string {
	return "fanout"
}


// Invoke actually performs the fanout
func (fanout *SomeFanout) Invoke(invocation *Invocation) (error, int) {
	invocation.ErrPipe.Drain()
	invocation.AutoHandleSignals()
	piChan, errChan := fanout.strategy(fanout, invocation)
	var wg sync.WaitGroup
	for {
		select {
		case pi, ok := <-piChan:
			if ok {
				wg.Add(1)
				go func (pi1 *PipelineInvocation) {
					pi1.Wait()
					wg.Done()
				}(pi)
			} else {
				break
			}
		case err := <-errChan:
			return err, 1
		}
	}
	wg.Wait()
	return nil, 0
}

// Factory for *SomeFanout
func NewFanout() *SomeFanout {
	fanout := new(SomeFanout)
	fanout.strategy = FanoutByLineStrategy
	return fanout
}

// Factory for *SomeFanout
func FanoutByLine(args ...Pipable) *SomeFanout {
	fanout := NewFanout()
	fanout.pipeline = NewPipeline(args...)
	return fanout
}

