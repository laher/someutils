package some

import (
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"github.com/laher/wget-go/wget"
	"io"
)

func init() {
	someutils.RegisterSome(func() someutils.SomeUtil { return NewWget() })
}

// SomeWget represents and performs a `wget` invocation
type SomeWget struct {
	// TODO: add members here
	args []string
}

// Name() returns the name of the util
func (w *SomeWget) Name() string {
	return "wget"
}

// TODO: add validation here

// ParseFlags parses flags from a commandline []string
func (w *SomeWget) ParseFlags(call []string, errWriter io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("wget", "[options] [args...]", someutils.VERSION)
	flagSet.SetOutput(errWriter)

	// TODO add flags here

	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errWriter, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}

	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	w.args = flagSet.Args()
	return nil
}

// Exec actually performs the wget
func (w *SomeWget) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	//TODO do something here!
	return wget.Wget(w.args)
}

// Factory for *SomeWget
func NewWget() *SomeWget {
	return new(SomeWget)
}

// Fluent factory for *SomeWget
func Wget(args ...string) *SomeWget {
	w := NewWget()
	w.args = args
	return w
}

// CLI invocation for *SomeWget
func WgetCli(call []string) error {
	w := NewWget()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := w.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return w.Exec(inPipe, outPipe, errPipe)
}
