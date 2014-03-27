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
func StdPipes() (io.Reader, io.Writer, io.Writer) {
	return os.Stdin, os.Stdout, os.Stderr
}

type Pipeline struct {
	InPipe  io.Reader
	OutPipe io.Writer
	ErrPipe io.Writer
	ErrInPipe io.Reader
}

func NewStdPipeline() *Pipeline {
	br := new(bytes.Reader)
	return &Pipeline{os.Stdin, os.Stdout, os.Stderr, br}
}
func NewPipeline(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) *Pipeline {
	br := new(bytes.Reader)
	return &Pipeline{inPipe, outPipe, errPipe, br}
}
func NewPipelineFromString(input string) (*Pipeline, *bytes.Buffer, *bytes.Buffer) {
	var outPipe bytes.Buffer
	var errPipe bytes.Buffer
	br := new(bytes.Reader)
	return &Pipeline{strings.NewReader(input), &outPipe, &errPipe, br}, &outPipe, &errPipe
}

func runAsync(pipable Pipable, inPipe io.Reader, outPipe io.Writer, errOutPipe io.Writer, errInPipe io.Reader, closers []io.Closer, e chan error) {
	_, willRedirectErrIn := pipable.(WillRedirectErrIn)
	if !willRedirectErrIn {
		go  func() {
			j, err := io.Copy(errOutPipe, errInPipe)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error! copying errInPipe to errOutPipe", err)
			}
			fmt.Fprintln(os.Stderr, "Finished copying errInPipe to errOutPipe", j)
		}()
	}
	go func() {
	e <- runSync(pipable, inPipe, outPipe, errOutPipe, errInPipe, closers)
	}()
}
func runSync(pipable Pipable, inPipe io.Reader, outPipe io.Writer, errOutPipe io.Writer, errInPipe io.Reader, closers []io.Closer) error {
	//_, willRedirectErrIn := pipable.(WillRedirectErrIn)
		//var locErrOutPipe io.Writer
	errSent := false
	/*
	var myInPipe io.Reader
	//var myOutPipe io.Writer
	//var myErrOutPipe io.Writer
	//bb, isBb := errOutPipe.(*bytes.Buffer)
	//if isBb {
		//prevErrOutput := bb.String()
		//r := strings.NewReader(prevErrOutput)
		//println("prev errOut: "+prevErrOutput)
	//}

	//myErrOutPipe := new(bytes.Buffer)
	if willRedirectErrIn {
		println("Running a WillRedirectErrIn")
		/*
		bb, isBb := (*errOutPipe).(*bytes.Buffer)
		if isBb {
			//prevErrOutput := bb.String()
			//r := strings.NewReader(prevErrOutput)
			//println("prev errOut: "+prevErrOutput)
			//bb.Truncate(0)
			//redirector.SetErrIn(r)
			//bb := bb.Read(0)
			//myInPipe = bytes.NewReader(bb.Bytes())
			//bb.Reset()
			// use a new writer
		} else {
			//Aargh.
			println("ERRRO")
			fmt.Fprintln(os.Stderr, "errPipe Not readable!")
		}
		* /
		myInPipe = errInPipe
	//	myOutPipe = outPipe
	} else {
		myInPipe = inPipe
	//	myOutPipe = outPipe
		//myErrOutPipe = errOutPipe
/*		go func() {
			i, err := io.Copy(errOutPipe, myErrOutPipe)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error! copying errInPipe to errOutPipe", e)
			}
			fmt.Fprintln(os.Stderr, "Finished copying errInPipe to errOutPipe", i)
		}()
		* /
	}
	*/
	if !errSent {
		err := pipable.Exec(inPipe, outPipe, errOutPipe)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "Ran pipable.Exec")
		/*
		bb, isBb := errOutPipe.(*bytes.Buffer)
		if isBb {
			println("ErrOut: ", bb.String())
		}
*/
	} else {
		//TODO show this has not run
		
		fmt.Fprintln(os.Stderr, "Could not run Exec")
	}
	var err error
	for _, closer := range closers {
		err = closer.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Close error ", err)
			if !errSent {
		//		return err
			}
		}
	}
	if !errSent {
		//e <- nil
		return err
	}
	return nil
}

func (p *Pipeline) Pipe(pipables ...Pipable) chan error {
	e := make(chan error)
	var previousReader *io.ReadCloser
	var previousErrReader *io.ReadCloser
	//fmt.Printf("%+v\n", pipables)
	for i, pipable := range pipables {
		//println(pipable)
		var w io.WriteCloser
		var r io.ReadCloser
		var wErr io.WriteCloser
		var rErr io.ReadCloser
		var locInPipe io.Reader
		var locOutPipe io.Writer
		var locErrInPipe io.Reader
		var locErrOutPipe io.Writer
		closers := []io.Closer{}
		if i == 0 {
			locInPipe = p.InPipe
			locErrInPipe = p.ErrInPipe
		} else {
			locInPipe = *previousReader
			locErrInPipe = *previousErrReader
		}
		if i == len(pipables)-1 {
			locOutPipe = p.OutPipe
			locErrOutPipe = p.ErrPipe
		} else {
			r, w = io.Pipe()
			locOutPipe = w
			closers = append(closers, w)

			rErr, wErr = io.Pipe()
			locErrOutPipe = wErr
			closers = append(closers, wErr)
		}
		//locErrOutPipe = &p.ErrPipe

/*
		bb, isBb := p.ErrPipe.(*bytes.Buffer)
		if isBb {
			//bytes of errPipe
			locErrInPipe = bytes.NewReader(bb.Bytes())
		} else {
			//nothing
			locErrInPipe = new(bytes.Buffer)
		}
		*/
		//go  func() {
			
			//println(pipable)
			runAsync(pipable, locInPipe, locOutPipe, locErrOutPipe, locErrInPipe, closers, e)

		//}()
		previousReader = &r
		previousErrReader = &rErr
	}
	return e
}

type WillRedirectErrIn interface {
	SetErrIn(errInPipe io.Reader)

}

func (p *Pipeline) PipeAndWait(timeoutSec time.Duration, pipables ...Pipable) (bool, []error) {
	e := p.Pipe(pipables...)
	return CollectErrors(e, len(pipables), timeoutSec)
}

func CollectErrors(e chan error, count int, timeoutSec time.Duration) (bool, []error) {
	errs := []error{}
	ok := true
	for i := 0; i < count; i++ {
		select {
		case <-time.After(timeoutSec*time.Second):
			errs = append(errs, errors.New("Timeout!"))
			return false, errs
		case err := <-e:
			if err != nil {
				ok = false
			}
			errs = append(errs, err)
		}
	}
	return ok, errs
}

type LineProcessorFunc func(io.Reader, io.Writer, io.Writer, []byte) error

func LineProcessor(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer, fu LineProcessorFunc) error {
	reader := bufio.NewReader(inPipe)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = fu(inPipe, outPipe, errPipe, line)
		if err != nil {
			return err
		}
	}
}
