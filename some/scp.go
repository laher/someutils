//+build !freebsd,!netbsd,!openbsd,!plan9

package some

import (
	"fmt"
	"github.com/laher/scp-go/scp"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
)

func init() {
	someutils.RegisterPipable(func() someutils.PipableCliUtil { return NewScp() })
}

// SomeScp represents and performs a `scp` invocation
type SomeScp struct {
	// TODO: add members here
	args []string
}

// Name() returns the name of the util
func (s *SomeScp) Name() string {
	return "scp"
}

// TODO: add validation here

// ParseFlags parses flags from a commandline []string
func (s *SomeScp) ParseFlags(call []string, errPipe io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("scp", "[options] [args...]", someutils.VERSION)
	flagSet.SetOutput(errPipe)

	// TODO add flags here

	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errPipe, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}

	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	s.args = flagSet.Args()
	// TODO: validate and process flagSet.Args()
	return nil
}

// Exec actually performs the scp
func (s *SomeScp) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	//TODO do something here!
	return scp.Scp(s.args)
}

// Factory for *SomeScp
func NewScp() *SomeScp {
	return new(SomeScp)
}

// Factory for *SomeScp
func Scp(args ...string) *SomeScp {
	s := NewScp()
	s.args = args
	return s
}

// CLI invocation for *SomeScp
func ScpCli(call []string) error {
	s := NewScp()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := s.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return s.Exec(inPipe, outPipe, errPipe)
}
