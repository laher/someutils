package some

import (
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"time"
)

func init() {
	someutils.RegisterPipable(func() someutils.PipableCliUtil { return NewTouch() })
}

// SomeTouch represents and performs a `touch` invocation
type SomeTouch struct {
	// TODO: add members here
	args []string
}

// Name() returns the name of the util
func (touch *SomeTouch) Name() string {
	return "touch"
}

// TODO: add validation here

// ParseFlags parses flags from a commandline []string
func (touch *SomeTouch) ParseFlags(call []string, errWriter io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("touch", "[options] [files...]", someutils.VERSION)
	flagSet.SetOutput(errWriter)

	// TODO add flags here

	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errWriter, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}

	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	args := flagSet.Args()
	if len(args) < 1 {
		return errors.New("Not enough args given")
	}
	touch.args = args

	return nil
}

// Exec actually performs the touch
func (touch *SomeTouch) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	//TODO do something here!
	for _, filename := range touch.args {
		err := touchFile(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func touchFile(filename string) error {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(filename)
			if err != nil {
				return err
			}
			return file.Close()
		} else {
			return err
		}
	} else {
		//set access times
		os.Chtimes(filename, time.Now(), time.Now())
	}
	return nil
}

// Factory for *SomeTouch
func NewTouch() *SomeTouch {
	return new(SomeTouch)
}

// Fluent factory for *SomeTouch
func Touch(args ...string) *SomeTouch {
	touch := NewTouch()
	touch.args = args
	return touch
}

// CLI invocation for *SomeTouch
func TouchCli(call []string) error {
	touch := NewTouch()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := touch.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return touch.Exec(inPipe, outPipe, errPipe)
}
