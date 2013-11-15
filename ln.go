package someutils

import (
	"errors"
	"github.com/laher/uggo"
)

const (
	LN_VERSION = "0.2.0"
)

type LnOptions struct {
	IsForce    bool
	IsSymbolic bool
}

func init() {
	Register(Util{
		"ln",
		Ln})
}

func Ln(call []string) error {
	options := LnOptions{}
	flagSet := uggo.NewFlagSetDefault("ln", "[options] TARGET LINK_NAME", LN_VERSION)
	flagSet.BoolVar(&options.IsSymbolic, "s", false, "Symbolic")
	flagSet.BoolVar(&options.IsForce, "f", false, "Force")

	e := flagSet.Parse(call[1:])
	if e != nil {
		println("Error parsing flags")
		return e
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	args := flagSet.Args()
	if len(args) < 2 {
		flagSet.Usage()
		return errors.New("Not enough args!")
	}
	target := args[0]
	linkName := args[1]
	return makeLink(target, linkName, options)

}
