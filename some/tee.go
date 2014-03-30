package some

import (
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
)

func init() {
	someutils.RegisterPipable(func() someutils.PipableCliUtil { return NewTee() })
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
func (tee *SomeTee) ParseFlags(call []string, errPipe io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("tee", "[OPTION]... [FILE]...", someutils.VERSION)
	flagSet.SetOutput(errPipe)
	flagSet.AliasedBoolVar(&tee.isAppend, []string{"a", "append"}, false, "Append instead of overwrite")

	// TODO add flags here

	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errPipe, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}

	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	tee.args = flagSet.Args()
	return nil
}

// Exec actually performs the tee
func (tee *SomeTee) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	flag := os.O_CREATE
	if tee.isAppend {
		flag = flag | os.O_APPEND
	}
	writeables := uggo.ToWriteableOpeners(tee.args, flag, 0666)
	files, err := uggo.OpenAll(writeables)
	if err != nil {
		return err
	}
	writers := []io.Writer{outPipe}
	for _, file := range files {
		writers = append(writers, file)
	}
	multiwriter := io.MultiWriter(writers...)
	_, err = io.Copy(multiwriter, inPipe)
	if err != nil {
		return err
	}
	for _, file := range files {
		err = file.Close()
		if err != nil {
			return err
		}
	}
	return nil

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
func TeeCli(call []string) error {
	tee := NewTee()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := tee.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return tee.Exec(inPipe, outPipe, errPipe)
}
