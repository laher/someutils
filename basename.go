package someutils

import (
	"github.com/laher/uggo"
	"path"
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
	base := path.Base(flagSet.Args()[0])
	println(base)
	return nil
}
