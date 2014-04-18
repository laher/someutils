package someutils

import (
	"fmt"
	"io"
	"os"
)

type Pipe struct {
	In  io.Reader
	Out io.Writer
}

func (p *Pipe) Drain() {
	go autoPipe(p.Out, p.In)
}

func (p *Pipe) CloseIfClosers() {
	closer, ok := p.In.(io.Closer)
	if ok {
		pipereader, ok := p.In.(*io.PipeReader)
		if ok {
			err := pipereader.CloseWithError(io.EOF)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not close inPipe (", err, ")")
			}
		} else {
			err := closer.Close()
			//never mind
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not close inPipe (", err, ")")
			}
		}
	} else {
		//		fmt.Fprintln(os.Stderr, "inPipe not a Closer")
	}
	closer, ok = p.Out.(io.Closer)
	if ok {
		pipewriter, ok := p.Out.(*io.PipeWriter)
		if ok {
			err := pipewriter.CloseWithError(io.EOF)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not close outPipe (", err, ")")
			}
		} else {
			err := closer.Close()
			//never mind
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not close outPipe (", err, ")")
			}
		}
	} else {
		//		fmt.Fprintln(os.Stderr, "outPipe not a Closer")
	}

}
