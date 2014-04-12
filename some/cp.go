package some

import (
	"errors"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
)

func init() {
	someutils.RegisterPipable(func() someutils.NamedPipable { return new(SomeCp) })
}

// SomeCp represents and performs a `cp` invocation
type SomeCp struct {
	// TODO: add members here
	IsRecursive bool
	SrcGlobs    []string
	Dest        string
}

// Name() returns the name of the util
func (cp *SomeCp) Name() string {
	return "cp"
}

// ParseFlags parses flags from a commandline []string
func (cp *SomeCp) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("cp", "[options] [src...] [dest]", someutils.VERSION)
	flagSet.AliasedBoolVar(&cp.IsRecursive, []string{"R", "r", "recursive"}, false, "Recurse into directories")
	flagSet.SetOutput(errPipe)

	// TODO add flags here
	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}
	args := flagSet.Args()
	if len(args) < 2 {
		flagSet.Usage()
		return errors.New("Not enough args"), 1
	}

	cp.SrcGlobs = args[0 : len(args)-1]
	cp.Dest = args[len(args)-1]

	return nil, 0
}

// Exec actually performs the cp
func (cp *SomeCp) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.AutoPipeErrInOut()
	invocation.AutoHandleSignals()
	for _, srcGlob := range cp.SrcGlobs {
		srces, err := filepath.Glob(srcGlob)
		if err != nil {
			return err, 1
		}
		for _, src := range srces {
			err = copyFile(src, cp.Dest, cp)
			if err != nil {
				return err, 1
			}
		}
	}
	return nil, 0
}

func copyFile(src, dest string, cp *SomeCp) error {
	//println("copy "+src+" to "+dest)

	srcFile, err := os.Open(src)
	defer srcFile.Close()
	if err != nil {
		return err
	}
	sinf, err := srcFile.Stat()
	if err != nil {
		return err
	}
	if sinf.IsDir() && !cp.IsRecursive {
		return errors.New("Omitting directory " + src)
	}

	//check if destination given is full filename or its (existing) parent dir
	var destFull string
	dinf, err := os.Stat(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} else {
			//doesnt exist yet. New file/dir
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
	//println("copy "+src+" to "+destFull)

	var destExists bool
	dinf, err = os.Stat(destFull)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} else {
			//doesnt exist. New file/dir
			destExists = false
		}
	} else {
		destExists = true
		if sinf.IsDir() && !dinf.IsDir() {
			return errors.New("destination is an existing non-directory")
		}
	}

	if sinf.IsDir() {
		//println("copying dir")
		if !destExists {
			//println("mkdir")
			err = os.Mkdir(destFull, sinf.Mode())
			if err != nil {
				return err
			}
		} else {
			//continue
		}
		contents, err := srcFile.Readdir(0)
		if err != nil {
			return err
		}
		err = srcFile.Close()
		if err != nil {
			return err
		}
		for _, fi := range contents {
			copyFile(filepath.Join(src, fi.Name()), destFull, cp)
		}
	} else {
		flags := os.O_WRONLY
		if !destExists {
			flags = flags + os.O_CREATE
		} else {
			flags = flags + os.O_TRUNC
		}
		destFile, err := os.OpenFile(destFull, flags, sinf.Mode())
		defer destFile.Close()
		if err != nil {
			return err
		}
		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}
		err = destFile.Close()
		if err != nil {
			return err
		}
		err = srcFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// Factory for *SomeCp
func NewCp() *SomeCp {
	return new(SomeCp)
}

// Factory for *SomeCp
func Cp(args ...string) *SomeCp {
	cp := NewCp()
	cp.SrcGlobs = args[0 : len(args)-1]
	cp.Dest = args[len(args)-1]
	return cp
}

// CLI invocation for *SomeCp
func CpCli(call []string) (error, int) {
	util := new(SomeCp)
	return someutils.StdInvoke((util), call)
}
