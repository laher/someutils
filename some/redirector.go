package some

import (
	"io"
	"io/ioutil"
	"os"
)

// PipeRedirector represents and performs a redirection between one Execable and another
// Note that PipeRedirector is an Execable but not a CLI util
type PipeRedirector struct {
	isRedirectFromErrPipe bool
	isRedirectToErrPipe   bool
	isRedirectToNull      bool
	isAppend              bool
	Filename              string
	errInPipe             io.Reader
}

// ErrPipeRedirector redirects the previous Execable's 'err' pipe
type ErrPipeRedirector struct {
	PipeRedirector
}

//redirects from errIn
func (redirector *ErrPipeRedirector) SetErrIn(errInPipe io.Reader) {
	redirector.errInPipe = errInPipe
}

// Exec actually performs the redirection
func (redirector *PipeRedirector) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	if redirector.Filename != "" {
		var fo io.WriteCloser
		flag := os.O_CREATE | os.O_WRONLY
		if redirector.isAppend {
			flag = flag | os.O_APPEND
		}
		fo, err := os.OpenFile(redirector.Filename, flag, 0777)
		if err != nil {
			return err
		}
		defer fo.Close()
		var rdr io.Reader
		if redirector.isRedirectFromErrPipe {
			rdr = redirector.errInPipe
		} else {
			rdr = inPipe
		}
		_, err = io.Copy(fo, rdr)
		if err != nil {
			return err
		}
		err = fo.Close()
		return err
	} else {
		var rdr io.Reader
		var writer io.Writer
		if redirector.isRedirectFromErrPipe {
			rdr = redirector.errInPipe
		} else {
			rdr = inPipe
		}

		if redirector.isRedirectToErrPipe {
			writer = errPipe
		} else if redirector.isRedirectToNull {
			writer = ioutil.Discard
		} else {
			writer = outPipe
		}

		closer, isCloser := writer.(io.Closer)
		if isCloser {
			// just incase
			defer closer.Close()
		}

		_, err := io.Copy(writer, rdr)
		if err != nil {
			return err
		}

		if isCloser {
			return closer.Close()
		}
		return nil
	}
}

// Factory for *PipeRedirector
func NewPipeRedirector() *PipeRedirector {
	return new(PipeRedirector)
}

// Factory for redirecting 'out' pipe to a file
func OutTo(filename string) *PipeRedirector {
	redirector := NewPipeRedirector()
	redirector.Filename = filename
	return redirector
}

// Factory for redirecting 'err' pipe to a file
func ErrTo(filename string) *PipeRedirector {
	redirector := NewPipeRedirector()
	redirector.Filename = filename
	redirector.isRedirectToErrPipe = true
	return redirector
}

// Factory for redirecting 'out' pipe to Null (nowhere)
func OutToNull() *ErrPipeRedirector {
	redirector := new(ErrPipeRedirector)
	redirector.isRedirectFromErrPipe = true
	redirector.isRedirectToNull = true
	return redirector
}

// Factory for redirecting 'err' pipe to Null (nowhere)
func ErrToNull() *ErrPipeRedirector {
	redirector := new(ErrPipeRedirector)
	redirector.isRedirectFromErrPipe = true
	redirector.isRedirectToNull = true
	return redirector
}

// Factory for redirecting 'err' pipe to 'out' pipe
func ErrToOut() *ErrPipeRedirector {
	redirector := new(ErrPipeRedirector)
	redirector.isRedirectFromErrPipe = true
	return redirector
}

// Factory for redirecting 'out' pipe to 'err' pipe
func OutToErr() *PipeRedirector {
	redirector := NewPipeRedirector()
	redirector.isRedirectToErrPipe = true
	return redirector
}
