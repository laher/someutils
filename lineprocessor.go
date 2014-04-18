package someutils

import (
	"bufio"
	"io"
)

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
