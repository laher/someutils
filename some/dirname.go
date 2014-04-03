package some

import (
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"path"
)

func init() {
	someutils.RegisterPipable(func() someutils.NamedPipable { return NewDirname() })
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
func (dirname *SomeDirname) ParseFlags(call []string, errPipe io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("dirname", "[options] NAME...", someutils.VERSION)
	flagSet.SetOutput(errPipe)

	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errPipe, "Flag error:  %v\n\n", err.Error())
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
func (dirname *SomeDirname) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	for _, f := range dirname.Filenames {
		dir := path.Dir(f)
		fmt.Fprintln(outPipe, dir)
	}
	return nil
}

// Factory for *SomeDirname
func NewDirname() *SomeDirname {
	return new(SomeDirname)
}

// Factory for *SomeDirname
func Dirname(args ...string) *SomeDirname {
	dirname := NewDirname()
	dirname.Filenames = args
	return dirname
}

// CLI invocation for *SomeDirname
func DirnameCli(call []string) error {
	dirname := NewDirname()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := dirname.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return dirname.Exec(inPipe, outPipe, errPipe)
}
