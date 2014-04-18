package someutils

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
}

/*
// ErrPipeRedirector redirects the previous Execable's 'err' pipe
type ErrPipeRedirector struct {
	PipeRedirector
}
*/

// Exec actually performs the redirection
func (redirector *PipeRedirector) Invoke(invocation *Invocation) (error, int) {
	if redirector.Filename != "" {
		var fo io.WriteCloser
		flag := os.O_CREATE | os.O_WRONLY
		if redirector.isAppend {
			flag = flag | os.O_APPEND
		}
		fo, err := os.OpenFile(redirector.Filename, flag, 0777)
		if err != nil {
			return err, 1
		}
		defer fo.Close()
		var rdr io.Reader
		if redirector.isRedirectFromErrPipe {
			rdr = invocation.ErrPipe.In
		} else {
			rdr = invocation.MainPipe.In
		}
		_, err = io.Copy(fo, rdr)
		if err == io.EOF || err == io.ErrClosedPipe {
			// OK
		} else if err != nil {
			return err, 1
		}
		err = fo.Close()
		if err != nil {
			return err, 1
		}
	} else {
		var rdr io.Reader
		var writer io.Writer
		if redirector.isRedirectFromErrPipe {
			rdr = invocation.ErrPipe.In
		} else {
			rdr = invocation.MainPipe.In
		}

		if redirector.isRedirectToErrPipe {
			writer = invocation.ErrPipe.Out
		} else if redirector.isRedirectToNull {
			//wrap Discard into an io.Pipe to ensure it is Closable
			r, w := io.Pipe()
			go io.Copy(ioutil.Discard, r)
			writer = w
		} else {
			writer = invocation.MainPipe.Out
		}

		closer, isCloser := writer.(io.Closer)
		if isCloser {
			// just incase
			defer closer.Close()
		}

		_, err := io.Copy(writer, rdr)
		if err == io.EOF || err == io.ErrClosedPipe {
			// OK
		} else if err != nil {
			return err, 1
		}

		if isCloser {
			err = closer.Close()
			if err != nil {
				return err, 1
			}
		}
	}
	return nil, 0
}

// Factory for redirecting 'out' pipe to a file
func OutTo(filename string) *PipeRedirector {
	redirector := new(PipeRedirector)
	redirector.Filename = filename
	return redirector
}

// Factory for redirecting 'err' pipe to a file
func ErrTo(filename string) *PipeRedirector {
	redirector := new(PipeRedirector)
	redirector.Filename = filename
	redirector.isRedirectToErrPipe = true
	return redirector
}

// Factory for redirecting 'out' pipe to Null (nowhere)
func OutToNull() *PipeRedirector {
	redirector := new(PipeRedirector)
	redirector.isRedirectToNull = true
	return redirector
}

// Factory for redirecting 'err' pipe to Null (nowhere)
func ErrToNull() *PipeRedirector {
	redirector := new(PipeRedirector)
	redirector.isRedirectFromErrPipe = true
	redirector.isRedirectToNull = true
	return redirector
}

// Factory for redirecting 'err' pipe to 'out' pipe
func ErrToOut() *PipeRedirector {
	redirector := new(PipeRedirector)
	redirector.isRedirectFromErrPipe = true
	return redirector
}

// Factory for redirecting 'out' pipe to 'err' pipe
func OutToErr() *PipeRedirector {
	redirector := new(PipeRedirector)
	redirector.isRedirectToErrPipe = true
	return redirector
}
