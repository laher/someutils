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
	someutils.RegisterPipable(func() someutils.CliPipable { return new(SomeHead) })
}

// SomeHead represents and performs a `head` invocation
type SomeHead struct {
	lines     int
	ch        byte
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
func (head *SomeHead) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.ErrPipe.Drain()
	invocation.AutoHandleSignals()
	//TODO do something here!
	if len(head.Filenames) > 0 {
		for _, fileName := range head.Filenames {
			file, err := os.Open(fileName)
			if err != nil {
				return err, 1
			}
			//err = headFile(file, head, invocation.MainPipe.Out)
			err = head.head(invocation.MainPipe.Out, file)
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
		err := head.head(invocation.MainPipe.Out, invocation.MainPipe.In)
		if err != nil {
			return err, 1
		}
	}
	return nil, 0
}

func (head *SomeHead) head(out io.Writer, in io.Reader) error {
	reader := bufio.NewReader(in)
	lineNo := 1
	ch := '\n' //should this be an option?
	for lineNo <= head.lines {
		text, err := reader.ReadBytes(byte(ch))
		if err != nil {
			return err
		}
		//text := scanner.Text()
		fmt.Fprintf(out, "%s", text) //, string(ch))
		lineNo++
	}
	/*err := scanner.Err()
	if err != nil {
		return err
	}
	*/
	return nil
}

// deprecated (use of bufio.Scanner)
func headFile(file io.Reader, head *SomeHead, out io.Writer) error {
	scanner := bufio.NewScanner(file)
	lineNo := 1
	for scanner.Scan() && lineNo <= head.lines {
		text := scanner.Text()
		fmt.Fprintf(out, "%s\n", text)
		lineNo++
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
func Head(lines int, args ...string) someutils.CliPipable {
	head := NewHead()
	head.lines = lines
	head.Filenames = args
	return (head)
}

// CLI invocation for *SomeHead
func HeadCli(call []string) (error, int) {
	util := new(SomeHead)
	return someutils.StdInvoke((util), call)
}
