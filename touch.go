package someutils

import (
	"errors"
	"flag"
	"os"
	"time"
)

func init() {
	Register(Util{
		"touch",
		Touch})
}

func Touch(call []string) error {
	flagSet := flag.NewFlagSet("unzip", flag.ContinueOnError)
	helpFlag := flagSet.Bool("help", false, "Show this help")
	err := flagSet.Parse(splitSingleHyphenOpts(call[1:]))
	if err != nil {
		return err
	}

	if *helpFlag {
		println("`zip` [options] [files...]")
		flagSet.PrintDefaults()
		return nil
	}
	args := flagSet.Args()
	if len(args) < 1 {
		return errors.New("Not enough args given")
	}
	for _, filename := range args {
		err = touch(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func touch(filename string) error {
		_, err := os.Stat(filename)
		if err != nil {
			if os.IsNotExist(err) {
				file, err := os.Create(filename)
				if err != nil {
					return err
				}
				return file.Close()
			} else {
				return err
			}
		} else {
			//set access times
			os.Chtimes(filename, time.Now(), time.Now())
		}
	return nil
}