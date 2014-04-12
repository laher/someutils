package some

import (
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func init() {
	someutils.RegisterPipable(func() someutils.CliPipable { return new(SomeWhich) })
}

// SomeWhich represents and performs a `which` invocation
type SomeWhich struct {
	all  bool
	args []string
}

// Name() returns the name of the util
func (which *SomeWhich) Name() string {
	return "which"
}

// ParseFlags parses flags from a commandline []string
func (which *SomeWhich) ParseFlags(call []string, errWriter io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("which", "[options] [args...]", someutils.VERSION)
	flagSet.SetOutput(errWriter)

	flagSet.BoolVar(&which.all, "a", false, "Print all matching executables in PATH, not just the first.")

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}

	which.args = flagSet.Args()
	return nil, 0
}

// Exec actually performs the which
func (which *SomeWhich) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.AutoPipeErrInOut()
	invocation.AutoHandleSignals()
	path := os.Getenv("PATH")
	if runtime.GOOS == "windows" {
		path = ".;" + path
	}
	pl := filepath.SplitList(path)
	for _, arg := range which.args {
		checkPathParts(arg, pl, which, invocation.OutPipe)
	}
	return nil, 0

}

func checkPathParts(arg string, pathParts []string, which *SomeWhich, outPipe io.Writer) {
	for _, pathPart := range pathParts {
		fi, err := os.Stat(pathPart)
		if err == nil {
			if fi.IsDir() {
				possibleExe := filepath.Join(pathPart, arg)
				if runtime.GOOS == "windows" {
					if !strings.HasSuffix(possibleExe, ".exe") {
						possibleExe += ".exe"
					}
				}
				_, err := os.Stat(possibleExe)
				if err != nil {
					//skip
				} else {
					abs, err := filepath.Abs(possibleExe)
					if err == nil {
						fmt.Fprintln(outPipe, abs)
					} else {
						//skip
					}
					if !which.all {
						return
					}
				}
			} else {
				//skip
			}
		} else {
			//skip
		}
	}
}

// Factory for *SomeWhich
func NewWhich() *SomeWhich {
	return new(SomeWhich)
}

// Factory for *SomeWhich
func Which(args ...string) *SomeWhich {
	which := NewWhich()
	which.args = args
	return which
}

// CLI invocation for *SomeWhich
func WhichCli(call []string) (error, int) {
	util := new(SomeWhich)
	return someutils.StdInvoke((util), call)
}
