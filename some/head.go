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
	someutils.RegisterSome(func() someutils.SomeUtil { return NewHead() })
}

// SomeHead represents and performs a `head` invocation
type SomeHead struct {
	lines int
	Filenames []string
}

// Name() returns the name of the util
func (head *SomeHead) Name() string {
	return "head"
}

// TODO: add validation here

// ParseFlags parses flags from a commandline []string
func (head *SomeHead) ParseFlags(call []string, errWriter io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("head", "[options] [args...]", someutils.VERSION)
	flagSet.SetOutput(errWriter)

	flagSet.AliasedIntVar(&head.lines, []string{"n", "lines"}, 10, "number of lines to print")
	
	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errWriter, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}

	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	//could be nil
	head.Filenames = flagSet.Args()
	return nil
}

// Exec actually performs the head
func (head *SomeHead) Exec(pipes someutils.Pipes) error {
	//TODO do something here!
	if len(head.Filenames)>0 {
		for _, fileName := range head.Filenames {
			file, err := os.Open(fileName)
			if err != nil {
				return err
			}
			err = headFile(file, head, pipes.Out())
			if err != nil {
				file.Close()
				return err
			}
			err = file.Close()
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		//stdin ..
		return headFile(pipes.In(), head, pipes.Out())
	}
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

// Fluent factory for *SomeHead
func Head(args ...string) *SomeHead {
	head := NewHead()
	head.Filenames = args
	return head
}

// CLI invocation for *SomeHead
func HeadCli(call []string) error {
	head := NewHead()
	pipes := someutils.StdPipes()
	err := head.ParseFlags(call, pipes.Err())
	if err != nil {
		return err
	}
	return head.Exec(pipes)
}
