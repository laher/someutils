package some

import (
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
)

func init() {
	someutils.RegisterSome(func() someutils.SomeUtil { return NewPwd() })
}

// SomePwd represents and performs a `pwd` invocation
type SomePwd struct {
	// TODO: add members here
}

// Name() returns the name of the util
func (pwd *SomePwd) Name() string {
	return "pwd"
}

// TODO: add validation here

// ParseFlags parses flags from a commandline []string
func (pwd *SomePwd) ParseFlags(call []string, errPipe io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("pwd", "", someutils.VERSION)
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

	// TODO: validate and process flagSet.Args()
	return nil
}

// Exec actually performs the pwd
func (pwd *SomePwd) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprintln(outPipe, wd)
	return nil
}

// Factory for *SomePwd
func NewPwd() *SomePwd {
	return new(SomePwd)
}

// Fluent factory for *SomePwd
func Pwd(args ...string) *SomePwd {
	pwd := NewPwd()
	return pwd
}

// CLI invocation for *SomePwd
func PwdCli(call []string) error {
	pwd := NewPwd()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := pwd.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return pwd.Exec(inPipe, outPipe, errPipe)
}
