package someutils

import (
	"flag"
	"os"
)

func init() {
	Register(Util{
		"pwd",
		Pwd})
}

func Pwd(call []string) error {

	flagSet := flag.NewFlagSet("pwd", flag.ContinueOnError)
	helpFlag := flagSet.Bool("help", false, "Show this help")

	err := flagSet.Parse(splitSingleHyphenOpts(call[1:]))
	if err != nil {
		return err
	}

	if *helpFlag {
		println("`pwd` [options] [files...]")
		flagSet.PrintDefaults()
		return nil
	}
	
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	println(wd)
	return nil
}