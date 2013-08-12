package someutils

import (
	"flag"
)

type CpOptions struct {
	recursive *bool
}

func Cp(call []string) error {
	options := CpOptions{}
	flagSet := flag.NewFlagSet("ls", flag.ContinueOnError)
	options.recursive = flagSet.Bool("r", false, "Recurse into directories")
	helpFlag := flagSet.Bool("help", false, "Show this help")

	e := flagSet.Parse(splitSingleHyphenOpts(call[1:]))
	if e != nil {
		return e
	}

	if *helpFlag {
		println("`cp` [options] [src] [dest]")
		flagSet.PrintDefaults()
		return nil
	}

	args := flagSet.Args()

	if len(args) != 2 {
		println("`cp` [options] [src] [dest]")
		flagSet.PrintDefaults()
		return nil
	}

	return nil
}
