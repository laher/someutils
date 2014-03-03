package someutils

import (
	"errors"
	"fmt"
	"github.com/laher/uggo"
	"os/exec"
	"strings"
)

func init() {
	Register(Util{
		"kill",
		Kill})
}

// very basic way to kill process on windows
func Kill(call []string) error {

	flagSet := uggo.NewFlagSetDefault("kill", "", VERSION)

	err := flagSet.Parse(call[1:])
	if err != nil {
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	args := flagSet.Args()

	if len(args) < 1 {
		return errors.New("Missing process name")
	}

	for i := 0; i < len(args); i++ {
		name := args[i]
		b, err := exec.Command("taskkill", "/f", "/im", name+"*").Output()
		if err != nil {
			return err
		}
		if len(strings.Split(string(b), "\n")) <= 0 {
			fmt.Printf("Can't kill %s \n", name)
		} else {
			fmt.Printf("Killed %s\n", name)
		}
	}

	return nil
}
