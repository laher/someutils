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
	someutils.RegisterPipable(func() someutils.NamedPipable { return NewWhich() })
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

// TODO: add validation here

// ParseFlags parses flags from a commandline []string
func (which *SomeWhich) ParseFlags(call []string, errWriter io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("which", "[options] [args...]", someutils.VERSION)
	flagSet.SetOutput(errWriter)

	flagSet.BoolVar(&which.all, "a", false, "Print all matching executables in PATH, not just the first.")

	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errWriter, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}

	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	which.args = flagSet.Args()
	return nil
}

// Exec actually performs the which
func (which *SomeWhich) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	path := os.Getenv("PATH")
	if runtime.GOOS == "windows" {
		path = ".;" + path
	}
	pl := filepath.SplitList(path)
	for _, arg := range which.args {
		checkPathParts(arg, pl, which, outPipe)
	}
	return nil

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
func WhichCli(call []string) error {
	which := NewWhich()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := which.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return which.Exec(inPipe, outPipe, errPipe)
}
