package someutils

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
)

type WhichOptions struct {
	all   *bool
}

func init() {
	Register(Util{
		"which",
		Which})
}

func Which(call []string) error {
	options := WhichOptions{}
	flagSet := flag.NewFlagSet("which", flag.ContinueOnError)
	options.all = flagSet.Bool("a", false, "Print all matching executables in PATH, not just the first.")
	
	helpFlag := flagSet.Bool("help", false, "Show this help")
	
	err := flagSet.Parse(splitSingleHyphenOpts(call[1:]))
	if err != nil {
		println("Error parsing flags")
		return err
	}

	if *helpFlag {
		println("`ls` [options] [dirs...]")
		flagSet.PrintDefaults()
		return nil
	}

	args := flagSet.Args()
	path := os.Getenv("PATH")
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
					possibleExe += ".exe"
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
					if !*options.all {
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
