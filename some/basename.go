package some

import (
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"path"
	"strings"
)

func init() {
	someutils.RegisterSome(func() someutils.SomeUtil { return NewBasename() })
}

// SomeBasename represents and performs a `basename` invocation
type SomeBasename struct {
	InputPath  string
	RelativeTo string
}

// Name() returns the name of the util
func (basename *SomeBasename) Name() string {
	return "basename"
}

// ParseFlags parses flags from a commandline []string
func (basename *SomeBasename) ParseFlags(call []string, errWriter io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("basename", "", someutils.VERSION)
	flagSet.SetOutput(errWriter)
	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errWriter, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}

	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	if len(flagSet.Args()) < 1 {
		return errors.New("Missing operand")
	}
	if len(flagSet.Args()) > 1 {
		basename.RelativeTo = flagSet.Args()[0]
		basename.InputPath = flagSet.Args()[1]
	} else {
		basename.InputPath = flagSet.Args()[0]
	}
	return nil
}

// Exec actually performs the basename
func (basename *SomeBasename) Exec(pipes someutils.Pipes) error {
	if basename.RelativeTo != "" {
		last := strings.LastIndex(basename.RelativeTo, basename.InputPath)
		base := basename.InputPath[:last]
		_, err := fmt.Fprintln(pipes.Out(), base)
		return err
	} else {
		_, err := fmt.Fprintln(pipes.Out(), path.Base(basename.InputPath))
		return err
	}
}

// Factory for *SomeBasename
func NewBasename() *SomeBasename {
	return new(SomeBasename)
}

// Fluent factory for *SomeBasename
func Basename(args ...string) *SomeBasename {
	basename := NewBasename()
	//basename.Xxx = args
	return basename
}

// CLI invocation for *SomeBasename
func BasenameCli(call []string) error {
	basename := NewBasename()
	pipes := someutils.StdPipes()
	err := basename.ParseFlags(call, pipes.Err())
	if err != nil {
		return err
	}
	return basename.Exec(pipes)
}
