package someutils

import (
	"errors"
	"fmt"
	"github.com/laher/uggo"
	"os"
	"path/filepath"
)

func init() {
	Register(Util{
		"mv",
		mv})
}

func mv(call []string) error {
	flagSet := uggo.NewFlagSetDefault("mv", "[options] [src...] [dest]", VERSION)

	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	args := flagSet.Args()

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: not enough arguments\n\n")
		flagSet.Usage()
		return errors.New("Not enough arguments")
	}

	srcGlobs := args[0 : len(args)-1]
	dest := args[len(args)-1]
	for _, srcGlob := range srcGlobs {
		srces, err := filepath.Glob(srcGlob)
		if err != nil {
			return err
		}
		if len(srces) < 1 {
			return errors.New(fmt.Sprintf("Source glob '%s' does not match any files\n", srcGlob))
		}

		for _, src := range srces {
			err = moveFile(src, dest)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error %v\n", err)
				return err
			}
		}
	}
	return nil
}

func moveFile(src, dest string) error {
	fmt.Printf("%s -> %s\n", src, dest)

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
