package some

import (
	"io"
	"os/exec"
	"github.com/laher/someutils"
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
func (exe *SomeExec) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) (error, int) {
	cmd := exec.Command(exe.args[0], exe.args[1:]...)
	cmd.Stdin = inPipe
	cmd.Stdout = outPipe
	cmd.Stderr = errPipe
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
func Exec(args ...string) someutils.NamedPipable {
	exe := new(SomeExec)
	exe.args = args
	return someutils.WrapNamed(exe)
}
