package someutils

import (
	"errors"
	"github.com/laher/uggo"
	"os"
	"time"
)

func init() {
	Register(Util{
		"touch",
		Touch})
}

func Touch(call []string) error {
	flagSet := uggo.NewFlagSetDefault("touch", "[options] [files...]", VERSION)
	err := flagSet.Parse(call[1:])
	if err != nil {
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
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
