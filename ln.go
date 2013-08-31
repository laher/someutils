package someutils

import (
	"errors"
	"flag"
)

type LnOptions struct {
	IsForce   *bool
	IsSymbolic  *bool
}


func init() {
	Register(Util{
		"ln",
		Ln})
}

func Ln(call []string) error {
	options := LnOptions{}
	flagSet := flag.NewFlagSet("ln", flag.ContinueOnError)
	options.IsSymbolic = flagSet.Bool("s", false, "Symbolic")
	options.IsForce = flagSet.Bool("f", false, "Force")
	helpFlag := flagSet.Bool("help", false, "Show this help")

	e := flagSet.Parse(splitSingleHyphenOpts(call[1:]))
	if e != nil {
		println("Error parsing flags")
		return e
	}

	if *helpFlag {
		println("`ln` [options] TARGET LINK_NAME")
		flagSet.PrintDefaults()
		return nil
	}

	args := flagSet.Args()
	if len(args) < 2 {
		return errors.New("Not enough args!")
	}
	target := args[0]
	linkName := args[1]
	return makeLink(target, linkName, options)
	
}

