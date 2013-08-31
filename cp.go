package someutils

import (
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
)

func init() {
	Register(Util{
		"cp",
		Cp})
}

type CpOptions struct {
	Recursive *bool
}

func Cp(call []string) error {
	options := CpOptions{}
	flagSet := flag.NewFlagSet("ls", flag.ContinueOnError)
	options.Recursive = flagSet.Bool("r", false, "Recurse into directories")
	helpFlag := flagSet.Bool("help", false, "Show this help")

	err := flagSet.Parse(splitSingleHyphenOpts(call[1:]))
	if err != nil {
		return err
	}

	if *helpFlag {
		println("`cp` [options] [src] [dest]")
		flagSet.PrintDefaults()
		return nil
	}

	args := flagSet.Args()

	if len(args) < 2 {
		println("`cp` [options] [src...] [dest]")
		flagSet.PrintDefaults()
		return nil
	}
	
	srcGlobs := args[0:len(args)-1]
	dest:= args[len(args)-1]
	//fmt.Printf("globs %v\n", srcGlobs)
	for _, srcGlob := range srcGlobs {
		srces, err := filepath.Glob(srcGlob)
		if err != nil {
			return err
		}
		//fmt.Printf(" %v\n", srces)
		for _, src := range srces {
			err = copyFile(src, dest, options)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func copyFile(src, dest string, options CpOptions) error {
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
	if sinf.IsDir() && !*options.Recursive {
		return errors.New("Omitting directory "+src)
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
	//println("copy "+src+" to "+destFull)

	var destExists bool
	dinf, err = os.Stat(destFull)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} else {
			//doesnt exist
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
			copyFile(filepath.Join(src, fi.Name()), destFull, options)
		}
	} else {
		destFile, err := os.OpenFile(destFull, os.O_CREATE, sinf.Mode())
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
