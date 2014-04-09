package some

import (
	"bufio"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
)

func init() {
	someutils.RegisterSimple(func() someutils.CliPipableSimple { return new(SomeHead) })
}

// SomeHead represents and performs a `head` invocation
type SomeHead struct {
	lines     int
	Filenames []string
}

// Name() returns the name of the util
func (head *SomeHead) Name() string {
	return "head"
}

// ParseFlags parses flags from a commandline []string
func (head *SomeHead) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("head", "[options] [args...]", someutils.VERSION)
	flagSet.SetOutput(errPipe)

	flagSet.AliasedIntVar(&head.lines, []string{"n", "lines"}, 10, "number of lines to print")

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}
	//could be nil
	head.Filenames = flagSet.Args()
	return nil, 0
}

// Exec actually performs the head
func (head *SomeHead) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) (error, int) {
	//TODO do something here!
	if len(head.Filenames) > 0 {
		for _, fileName := range head.Filenames {
			file, err := os.Open(fileName)
			if err != nil {
				return err, 1
			}
			err = headFile(file, head, outPipe)
			if err != nil {
				file.Close()
				return err, 1
			}
			err = file.Close()
			if err != nil {
				return err, 1
			}
		}
	} else {
		//stdin ..
		err := headFile(inPipe, head, outPipe)
		if err != nil {
			return err, 1
		}
	}
	return nil, 0
}

func headFile(file io.Reader, head *SomeHead, out io.Writer) error {
	scanner := bufio.NewScanner(file)
	line := 1
	for scanner.Scan() && line <= head.lines {
		text := scanner.Text()
		fmt.Fprintf(out, "%s\n", text)
		line++
	}
	err := scanner.Err()
	if err != nil {
		return err
	}
	return nil
}

// Factory for *SomeHead
func NewHead() *SomeHead {
	return new(SomeHead)
}

// Factory for *SomeHead
func Head(lines int, args ...string) someutils.NamedPipable {
	head := NewHead()
	head.lines = lines
	head.Filenames = args
	return someutils.WrapNamed(head)
}

// CLI invocation for *SomeHead
func HeadCli(call []string) (error, int) {
	util := new(SomeHead)
	return someutils.StdInvoke(someutils.WrapUtil(util), call)
}
