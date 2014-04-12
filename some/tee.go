package some

import (
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
)

func init() {
	someutils.RegisterPipable(func() someutils.CliPipable { return new(SomeTee) })
}

// SomeTee represents and performs a `tee` invocation
type SomeTee struct {
	isAppend bool
	flag     int
	args     []string
}

// Name() returns the name of the util
func (tee *SomeTee) Name() string {
	return "tee"
}

// TODO: add validation here

// ParseFlags parses flags from a commandline []string
func (tee *SomeTee) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("tee", "[OPTION]... [FILE]...", someutils.VERSION)
	flagSet.SetOutput(errPipe)
	flagSet.AliasedBoolVar(&tee.isAppend, []string{"a", "append"}, false, "Append instead of overwrite")

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}

	tee.args = flagSet.Args()
	return nil, 0
}

// Exec actually performs the tee
func (tee *SomeTee) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.AutoPipeErrInOut()
	invocation.AutoHandleSignals()
	flag := os.O_CREATE
	if tee.isAppend {
		flag = flag | os.O_APPEND
	}
	writeables := uggo.ToWriteableOpeners(tee.args, flag, 0666)
	files, err := uggo.OpenAll(writeables)
	if err != nil {
		return err, 1
	}
	writers := []io.Writer{invocation.OutPipe}
	for _, file := range files {
		writers = append(writers, file)
	}
	multiwriter := io.MultiWriter(writers...)
	_, err = io.Copy(multiwriter, invocation.InPipe)
	if err != nil {
		return err, 1
	}
	for _, file := range files {
		err = file.Close()
		if err != nil {
			return err, 1
		}
	}
	return nil, 0

}

// Factory for *SomeTee
func NewTee() *SomeTee {
	return new(SomeTee)
}

// Factory for *SomeTee
func Tee(args ...string) *SomeTee {
	tee := NewTee()
	tee.args = args
	return tee
}

// CLI invocation for *SomeTee
func TeeCli(call []string) (error, int) {

	util := new(SomeTee)
	return someutils.StdInvoke((util), call)
}
