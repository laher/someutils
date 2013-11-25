package someutils

import (
	"github.com/laher/uggo"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type WhichOptions struct {
	all bool
}

func init() {
	Register(Util{
		"which",
		Which})
}

func Which(call []string) error {
	options := WhichOptions{}
	flagSet := uggo.NewFlagSetDefault("which", "[-a] args", VERSION)
	flagSet.BoolVar(&options.all, "a", false, "Print all matching executables in PATH, not just the first.")

	err := flagSet.Parse(call[1:])
	if err != nil {
		println("Error parsing flags")
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	args := flagSet.Args()
	path := os.Getenv("PATH")
	if runtime.GOOS == "windows" {
		path = ".;" + path
	}
	pl := filepath.SplitList(path)
	for _, arg := range args {
		checkPathParts(arg, pl, options)
		/*
			if err != nil {
				return err
			}*/
	}
	return nil
}

func checkPathParts(arg string, pathParts []string, options WhichOptions) {
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
						println(abs)
					} else {
						//skip
					}
					if !options.all {
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
