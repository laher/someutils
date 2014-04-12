package some

import (
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"path"
)

func init() {
	someutils.RegisterPipable(func() someutils.CliPipable { return new(SomeDirname) })
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
func (dirname *SomeDirname) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("dirname", "[options] NAME...", someutils.VERSION)
	flagSet.SetOutput(errPipe)

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}
	dirname.Filenames = flagSet.Args()
	return nil, 0
}

// Exec actually performs the dirname
func (dirname *SomeDirname) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.AutoPipeErrInOut()
	invocation.AutoHandleSignals()
	for _, f := range dirname.Filenames {
		dir := path.Dir(f)
		fmt.Fprintln(invocation.OutPipe, dir)
	}
	return nil, 0
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
func DirnameCli(call []string) (error, int) {
	util := new(SomeDirname)
	return someutils.StdInvoke((util), call)
}
