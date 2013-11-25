package someutils

import (
	"fmt"
	"github.com/laher/uggo"
	"io/ioutil"
	"os"
	"path/filepath"
)

func init() {
	Register(Util{
		"rm",
		Rm})
}

type RmOptions struct {
	IsRecursive bool
}

func Rm(call []string) error {
	options := RmOptions{}
	flagSet := uggo.NewFlagSetDefault("rm", "[options] [files...]", VERSION)
	flagSet.BoolVar(&options.IsRecursive, "r", false, "Recurse into directories")

	e := flagSet.Parse(call[1:])
	if e != nil {
		return e
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	for _, fileGlob := range flagSet.Args() {
		files, err := filepath.Glob(fileGlob)
		if err != nil {
			return err
		}
		for _, file := range files {
			e := delete(file, options.IsRecursive)
			if e != nil {
				return e
			}
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
