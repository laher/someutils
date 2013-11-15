package someutils

import (
	"github.com/laher/uggo"
	"os"
	"path/filepath"
)

const (
	MV_VERSION = "0.2.0"
)

func init() {
	Register(Util{
		"mv",
		mv})
}

type MvOptions struct {
	IsHelp bool
}

func mv(call []string) error {
	options := MvOptions{}
	flagSet := uggo.NewFlagSetDefault("mv", "[options] [src...] [dest]", MV_VERSION)
	flagSet.BoolVar(&options.IsHelp, "help", false, "Show this help")

	err := flagSet.Parse(call[1:])
	if err != nil {
		return err
	}

	if options.IsHelp {
		println("`mv` [options] [src] [dest]")
		flagSet.PrintDefaults()
		return nil
	}

	args := flagSet.Args()

	if len(args) < 2 {
		println("`mv` [options] [src...] [dest]")
		flagSet.PrintDefaults()
		return nil
	}

	srcGlobs := args[0 : len(args)-1]
	dest := args[len(args)-1]
	//fmt.Printf("globs %v\n", srcGlobs)
	for _, srcGlob := range srcGlobs {
		srces, err := filepath.Glob(srcGlob)
		if err != nil {
			return err
		}
		//fmt.Printf(" %v\n", srces)
		for _, src := range srces {
			err = moveFile(src, dest, options)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func moveFile(src, dest string, options MvOptions) error {

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
