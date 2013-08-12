package someutils

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type RmOptions struct {
	recursive *bool
}

func Rm(call []string) error {
	options := RmOptions{}
	flagSet := flag.NewFlagSet("ls", flag.ContinueOnError)
	options.recursive = flagSet.Bool("r", false, "Recurse into directories")
	helpFlag := flagSet.Bool("help", false, "Show this help")

	e := flagSet.Parse(splitSingleHyphenOpts(call[1:]))
	if e != nil {
		return e
	}

	if *helpFlag {
		println("`rm` [options] [files...]")
		flagSet.PrintDefaults()
		return nil
	}
	for _, file := range flagSet.Args() {
		e := delete(file, *options.recursive)
		if e != nil {
			return e
		}
	}

	return nil
}

func delete(file string, recursive bool) error {
	fi, e := os.Stat(file)
	if e != nil {
		return e
	}
	if fi.IsDir() && recursive {
		e := deleteDir(file)
		if e != nil {
			return e
		}
	} else if fi.IsDir() {
		//do nothing
		return fmt.Errorf("'%s' is a directory. Use -r", file)
	}
	return os.Remove(file)
}

func deleteDir(dir string) error {
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		return e
	}
	for _, file := range files {
		e = delete(filepath.Join(dir, file.Name()), true)
		if e != nil {
			return e
		}
	}
	return nil
}
