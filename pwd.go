package someutils

import (
	"github.com/laher/uggo"
	"os"
)

const (
	PWD_VERSION = "0.2.0"
)

func init() {
	Register(Util{
		"pwd",
		Pwd})
}

func Pwd(call []string) error {

	flagSet := uggo.NewFlagSetDefault("pwd", "", PWD_VERSION)

	err := flagSet.Parse(call[1:])
	if err != nil {
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	println(wd)
	return nil
}
