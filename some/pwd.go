package some

import (
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
)

func init() {
	someutils.RegisterPipable(func() someutils.NamedPipable { return NewPwd() })
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
func (pwd *SomePwd) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("pwd", "", someutils.VERSION)
	flagSet.SetOutput(errPipe)

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}
	return nil, 0
}

// Exec actually performs the pwd
func (pwd *SomePwd) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) (error, int) {
	wd, err := os.Getwd()
	if err != nil {
		return err, 1
	}
	fmt.Fprintln(outPipe, wd)
	return nil, 0
}

// Factory for *SomePwd
func NewPwd() *SomePwd {
	return new(SomePwd)
}

// Factory for *SomePwd
func Pwd(args ...string) *SomePwd {
	pwd := NewPwd()
	return pwd
}

// CLI invocation for *SomePwd
func PwdCli(call []string) (error, int) {
	pwd := NewPwd()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err, code := pwd.ParseFlags(call, errPipe)
	if err != nil {
		return err, code
	}
	return pwd.Exec(inPipe, outPipe, errPipe)
}
