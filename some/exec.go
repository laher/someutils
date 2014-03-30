package some

import (
	"io"
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
func (exe *SomeExec) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	cmd := exec.Command(exe.args[0], exe.args[1:]...)
	cmd.Stdin = inPipe
	cmd.Stdout = outPipe
	cmd.Stderr = errPipe
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

// Factory for *SomeExec
func NewExec() *SomeExec {
	return new(SomeExec)
}

// Factory for *SomeExec
func Exec(args ...string) *SomeExec {
	exe := NewExec()
	exe.args = args
	return exe
}
