package someutils

import (
	"flag"
	"os"
)

type MvOptions struct {
	recursive *bool
}

func Mv(call []string) error {
	options := MvOptions{}
	flagSet := flag.NewFlagSet("ls", flag.ContinueOnError)
	options.recursive = flagSet.Bool("r", false, "Recurse into directories")
	helpFlag := flagSet.Bool("help", false, "Show this help")

	err := flagSet.Parse(splitSingleHyphenOpts(call[1:]))
	if err != nil {
		return err
	}

	if *helpFlag {
		println("`mv` [options] [src] [dest]")
		flagSet.PrintDefaults()
		return nil
	}

	args := flagSet.Args()

	if len(args) != 2 {
		println("`mv` [options] [src] [dest]")
		flagSet.PrintDefaults()
		return nil
	}
	err = os.Rename(args[0], args[1])
	return err
}
