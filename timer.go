package someutils

import (
	"errors"
	"fmt"
	"github.com/laher/uggo"
	"os"
	"os/exec"
	"time"
)

func init() {
	Register(Util{
	"timer",
	Timer})
}

// very basic way to Time execution on windows
func Timer(call []string) error {
	var cmd *exec.Cmd
	start := time.Now()
	fmt.Println()
	flagSet := uggo.NewFlagSetDefault("timer", "", VERSION)

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

	os.Args[0] = "/C"
	cmd = exec.Command("cmd", os.Args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Time: ", time.Since(start))
	fmt.Println()
	return nil
}
