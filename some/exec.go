package some

import (
	"github.com/laher/someutils"
	"os/exec"
)

// SomeExec represents and performs a `exec` invocation
type SomeExec struct {
	// TODO: add members here
	args []string
}

// Name() returns the name of the util
func (exe *SomeExec) Name() string {
	return "exec"
}

// Exec actually performs the exec
func (exe *SomeExec) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.ErrPipe.Drain()
	invocation.AutoHandleSignals()
	cmd := exec.Command(exe.args[0], exe.args[1:]...)
	cmd.Stdin = invocation.MainPipe.In
	cmd.Stdout = invocation.MainPipe.Out
	cmd.Stderr = invocation.ErrPipe.Out
	err := cmd.Start()
	if err != nil {
		return err, 1
	}
	err = cmd.Wait()
	exitSuccess := cmd.ProcessState.Success()
	var exitCode int
	if exitSuccess {
		exitCode = 0
	} else {
		// There should be a way to get the proper status on Unix.
		exitCode = 1
	}
	return err, exitCode
}

// Factory for *SomeExec
func Exec(args ...string) someutils.Pipable {
	exe := new(SomeExec)
	exe.args = args
	return (exe)
}
