package some

import (
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"path"
)

func init() {
	someutils.RegisterSome(func() someutils.SomeUtil { return NewDirname() })
}

// SomeDirname represents and performs a `dirname` invocation
type SomeDirname struct {
	Filenames []string
}

// Name() returns the name of the util
func (dirname *SomeDirname) Name() string {
	return "dirname"
}

// ParseFlags parses flags from a commandline []string
func (dirname *SomeDirname) ParseFlags(call []string, errWriter io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("dirname", "[options] NAME...", someutils.VERSION)
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
	dirname.Filenames = flagSet.Args()
	return nil
}

// Exec actually performs the dirname
func (dirname *SomeDirname) Exec(pipes someutils.Pipes) error {
	for _, f := range dirname.Filenames {
		dir := path.Dir(f)
		fmt.Fprintln(pipes.Out(), dir)
	}
	return nil
}

// Factory for *SomeDirname
func NewDirname() *SomeDirname {
	return new(SomeDirname)
}

// Fluent factory for *SomeDirname
func Dirname(args ...string) *SomeDirname {
	dirname := NewDirname()
	dirname.Filenames = args
	return dirname
}

// CLI invocation for *SomeDirname
func DirnameCli(call []string) error {
	dirname := NewDirname()
	pipes := someutils.StdPipes()
	err := dirname.ParseFlags(call, pipes.Err())
	if err != nil {
		return err
	}
	return dirname.Exec(pipes)
}
