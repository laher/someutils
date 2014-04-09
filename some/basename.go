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
	someutils.RegisterPipable(func() someutils.NamedPipable { return someutils.WrapNamed(new(SomeBasename)) })
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
func (basename *SomeBasename) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("basename", "", someutils.VERSION)
	flagSet.SetOutput(errPipe)
	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}
	if len(flagSet.Args()) < 1 {
		return errors.New("Missing operand"), 1
	}
	if len(flagSet.Args()) > 1 {
		basename.RelativeTo = flagSet.Args()[0]
		basename.InputPath = flagSet.Args()[1]
	} else {
		basename.InputPath = flagSet.Args()[0]
	}
	return nil, 1
}

// Exec actually performs the basename
func (basename *SomeBasename) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) (error, int) {
	if basename.RelativeTo != "" {
		last := strings.LastIndex(basename.RelativeTo, basename.InputPath)
		base := basename.InputPath[:last]
		_, err := fmt.Fprintln(outPipe, base)
		if err != nil {
			return err, 1
		}
	} else {
		_, err := fmt.Fprintln(outPipe, path.Base(basename.InputPath))
		if err != nil {
			return err, 1
		}
	}
	return nil, 0
}



// Factory for *SomeBasename
func Basename(args ...string) someutils.NamedPipable {
	basename := new(SomeBasename)
	//basename.Xxx = args
	return someutils.WrapNamed(basename)
}

// CLI invocation for *SomeBasename
func BasenameCli(call []string) (error, int) {
	basename := new(SomeBasename)
	return someutils.StdInvoke(someutils.WrapUtil(basename), call)
}
