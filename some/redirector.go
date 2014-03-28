package some

import (
	"io"
	"os"
)
//Note that PipeRedirector is just an Execable - not intended as a Util

// PipeRedirector represents and performs a `>` invocation
type PipeRedirector struct {
	isRedirectFromErrPipe bool
	isRedirectToErrPipe bool
	isAppend bool
	Filename string
	errInPipe io.Reader
}

type UpstreamErrPipeRedirector struct {
	PipeRedirector
}

//redirects from errIn
func (redirector *UpstreamErrPipeRedirector) SetErrIn(errInPipe io.Reader) {
	redirector.errInPipe = errInPipe
}

// Exec actually performs the redirector
func (redirector *PipeRedirector) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	if redirector.Filename != "" {
		var fo io.WriteCloser
		flag := os.O_CREATE|os.O_WRONLY
		if redirector.isAppend {
			flag = flag | os.O_APPEND
		}
		fo, err := os.OpenFile(redirector.Filename, flag, 0777)
		if err != nil {
			return err
		}
		defer fo.Close()
		var rdr io.Reader
		/*
		if redirector.isRedirectFromErrPipe {
			rdr = redirector.errInPipe
		} else {
			rdr = inPipe
		}
		*/
		rdr = inPipe
		_, err = io.Copy(fo, rdr)
		if err != nil {
			return err
		}
		err = fo.Close()
		return err
	} else {
		var rdr io.Reader
		var writer io.Writer
		/*
		if redirector.isRedirectFromErrPipe {
			rdr = redirector.errInPipe
		} else {
			rdr = inPipe
		}
		*/
		rdr = inPipe
		if redirector.isRedirectToErrPipe {
			writer = errPipe
		} else {
			writer = outPipe
		}

		_, err := io.Copy(writer, rdr)
		if err != nil {
			return err
		}
		/*
		closer, isCloser := writer.(io.Closer)
		if isCloser {
			return closer.Close()
		}
		return err
		*/
		return nil
	}
}

// Factory for *PipeRedirector
func NewPipeRedirector() *PipeRedirector {
	return new(PipeRedirector)
}

// Fluent factory for *PipeRedirector
func RedirectTo(filename string) *PipeRedirector {
	redirector := NewPipeRedirector()
	redirector.Filename = filename
	return redirector
}

// Fluent factory for *PipeRedirector
func RedirectErrTo(filename string) *PipeRedirector {
	redirector := NewPipeRedirector()
	redirector.Filename = filename
	redirector.isRedirectToErrPipe = true
	return redirector
}

func RedirectErrToOut() *UpstreamErrPipeRedirector {
	redirector := new(UpstreamErrPipeRedirector)
	redirector.isRedirectFromErrPipe = true
	return redirector
}
func RedirectOutToErr() *PipeRedirector {
	redirector := NewPipeRedirector()
	redirector.isRedirectToErrPipe = true
	return redirector
}

