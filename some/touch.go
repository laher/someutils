package some

import (
	"errors"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"time"
)

func init() {
	someutils.RegisterSimple(func() someutils.CliPipableSimple { return new(SomeTouch) })
}

// SomeTouch represents and performs a `touch` invocation
type SomeTouch struct {
	args []string
}

// Name() returns the name of the util
func (touch *SomeTouch) Name() string {
	return "touch"
}

// ParseFlags parses flags from a commandline []string
func (touch *SomeTouch) ParseFlags(call []string, errWriter io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("touch", "[options] [files...]", someutils.VERSION)
	flagSet.SetOutput(errWriter)

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}

	args := flagSet.Args()
	if len(args) < 1 {
		return errors.New("Not enough args given"), 1
	}
	touch.args = args
	return nil, 0
}

// Exec actually performs the touch
func (touch *SomeTouch) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) (error, int) {
	for _, filename := range touch.args {
		err := touchFile(filename)
		if err != nil {
			return err, 1
		}
	}
	return nil, 0
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

// Factory for *SomeTouch
func Touch(args ...string) *SomeTouch {
	touch := NewTouch()
	touch.args = args
	return touch
}

// CLI invocation for *SomeTouch
func TouchCli(call []string) (error, int) {

	util := new(SomeTouch)
	return someutils.StdInvoke(someutils.WrapUtil(util), call)
}
