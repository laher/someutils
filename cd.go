package someutils

import (
	"fmt"
	"github.com/laher/uggo"
	"io"
	"os"
)

// SomeCd represents and performs a `cd` invocation
type SomeCd struct {
	destDir string
}

// Name() returns the name of the util
func (cd *SomeCd) Name() string {
	return "cd"
}

// ParseFlags parses flags from a commandline []string
func (cd *SomeCd) ParseFlags(call []string, errPipe io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("cd", "[options] [args...]", VERSION)
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

// Exec actually performs the cd
func (cd *SomeCd) Invoke(invocation *Invocation) (error, int) {
	invocation.ErrPipe.Drain()
	invocation.AutoHandleSignals()
	err := os.Chdir(cd.destDir)
	if err != nil {
		return err, 1
	} else {
		return err, 0
	}
}

// Factory for *SomeCd
func NewCd() *SomeCd {
	return new(SomeCd)
}

// Factory for *SomeCd
func Cd(destDir string) *SomeCd {
	cd := NewCd()
	cd.destDir = destDir
	return cd
}
