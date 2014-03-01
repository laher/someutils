package someutils

import (
	"github.com/laher/uggo"
	"path"
)

func init() {
	Register(Util{
		"dirname",
		Dirname})
}

//very basic for the moment. No removal of suffix
func Dirname(call []string) error {

	flagSet := uggo.NewFlagSetDefault("dirname", "", VERSION)

	err := flagSet.Parse(call[1:])
	if err != nil {
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	for _, f := range flagSet.Args() {
		dir := path.Dir(f)
		println(dir)
	}
	return nil
}
