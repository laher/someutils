package someutils

import (
	"errors"
	"github.com/laher/uggo"
	"path"
	"strings"
)

func init() {
	Register(Util{
		"basename",
		Basename})
}

//very basic for the moment. No removal of suffix
func Basename(call []string) error {

	flagSet := uggo.NewFlagSetDefault("basename", "", VERSION)

	err := flagSet.Parse(call[1:])
	if err != nil {
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	if len(flagSet.Args()) < 1 {
		return errors.New("Missing operand")
	}
	base := path.Base(flagSet.Args()[0])
	if len(flagSet.Args()) > 1 {
		if strings.HasSuffix(base, flagSet.Args()[1]) {
			last := strings.LastIndex(base, flagSet.Args()[1])
			base = base[:last]
		}
	}
	println(base)
	return nil
}
