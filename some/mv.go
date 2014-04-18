package some

import (
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
)

func init() {
	someutils.RegisterPipable(func() someutils.CliPipable { return new(SomeMv) })
}

// SomeMv represents and performs a `mv` invocation
type SomeMv struct {
	srcGlobs []string
	dest     string
}

// Name() returns the name of the util
func (mv *SomeMv) Name() string {
	return "mv"
}

// ParseFlags parses flags from a commandline []string
func (mv *SomeMv) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("mv", "[options] [src...] [dest]", someutils.VERSION)
	flagSet.SetOutput(errPipe)

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}

	args := flagSet.Args()

	if len(args) < 2 {
		fmt.Fprintf(errPipe, "Error: not enough arguments\n\n")
		flagSet.Usage()
		return errors.New("Not enough arguments"), 1
	}

	mv.srcGlobs = args[0 : len(args)-1]
	mv.dest = args[len(args)-1]

	return nil, 0
}

// Exec actually performs the mv
func (mv *SomeMv) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.ErrPipe.Drain()
	invocation.AutoHandleSignals()
	for _, srcGlob := range mv.srcGlobs {
		srces, err := filepath.Glob(srcGlob)
		if err != nil {
			return err, 1
		}
		if len(srces) < 1 {
			return errors.New(fmt.Sprintf("Source glob '%s' does not match any files\n", srcGlob)), 1
		}

		for _, src := range srces {
			err = moveFile(src, mv.dest)
			if err != nil {
				fmt.Fprintf(invocation.ErrPipe.Out, "Error %v\n", err)
				return err, 1
			}
		}
	}
	return nil, 0

}

func moveFile(src, dest string) error {
	//fmt.Printf("%s -> %s\n", src, dest)

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	sinf, err := srcFile.Stat()
	if err != nil {
		return err
	}
	err = srcFile.Close()
	if err != nil {
		return err
	}

	//check if destination given is full filename or its (existing) parent dir
	var destFull string
	dinf, err := os.Stat(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} else {
			//doesnt exist
			destFull = dest
		}
	} else {
		if dinf.IsDir() {
			//copy file name
			destFull = filepath.Join(dest, sinf.Name())
		} else {
			destFull = dest
		}
	}
	err = os.Rename(src, destFull)
	return err
}

// Factory for *SomeMv
func NewMv() *SomeMv {
	return new(SomeMv)
}

// Factory for *SomeMv
func Mv(args ...string) *SomeMv {
	mv := NewMv()
	mv.srcGlobs = args[0 : len(args)-1]
	mv.dest = args[len(args)-1]
	return mv
}

// CLI invocation for *SomeMv
func MvCli(call []string) (error, int) {

	util := new(SomeMv)
	return someutils.StdInvoke((util), call)
}
