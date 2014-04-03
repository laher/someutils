package some

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
)

func init() {
	someutils.RegisterPipable(func () someutils.NamedPipable { return NewXargs() })
}

// SomeXargs represents and performs a `xargs` invocation
type SomeXargs struct {
	utilFactory someutils.PipableCliUtilFactory
	utilArgs []string
	maxProcesses int
}

// Name() returns the name of the util
func (xargs *SomeXargs) Name() string {
	return "xargs"
}

// ParseFlags parses flags from a commandline []string
func (xargs *SomeXargs) ParseFlags(call []string, errPipe io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("xargs", "[options] [args...]", someutils.VERSION)
	flagSet.SetOutput(errPipe)
	// TODO multiple processes at once
	flagSet.AliasedIntVar(&xargs.maxProcesses, []string{"P", "max-procs"}, 1, "Maximum processes")
	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errPipe, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}

	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	args := flagSet.Args()
	if len(args) < 1 {
		return errors.New("No command specified")
	}
	if !someutils.Exists(args[0]) {
		return errors.New(" Command does not exist.")
	}
	xargs.utilFactory = someutils.GetPipableCliUtilFactory(args[0])
	xargs.utilArgs = args
	return nil
}

// Exec actually performs the xargs
func (xargs *SomeXargs) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	util := xargs.utilFactory()
	args := xargs.newArgset(util.Name())
	reader := bufio.NewReader(inPipe)
	cont := true
	count := 0
	maxCount := 5
	for cont {
		if count >= maxCount {
			count = 0
			//fmt.Fprintf(errPipe, "args for '%s': %v\n", util.Name(), args)
			err := util.ParseFlags(args, errPipe)
			if err != nil {
				return err
			}
			err = util.Exec(inPipe, outPipe, errPipe)
			if err != nil {
				return err
			}
		}
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			cont = false
		} else if err != nil {
			return err
		} else {
			args = append(args, string(line))
			if err != nil {
				return err
			}
			count++
		}
	}
	if count > 0 {
		//fmt.Fprintf(errPipe, "args for '%s': %v\n", util.Name(), args)
		err := util.ParseFlags(args, errPipe)
		if err != nil {
			return err
		}
		err = util.Exec(inPipe, outPipe, errPipe)
		return err
	}
	return nil
}

func (xargs *SomeXargs) newArgset(cmdName string) []string {
	args := []string{ cmdName }
	args = append(args, xargs.utilArgs...)
	return args
}

// Factory for *SomeXargs
func NewXargs() *SomeXargs {
	return new(SomeXargs)
}

// Factory for *SomeXargs
func Xargs(utilFactory someutils.PipableCliUtilFactory, args ...string) *SomeXargs {
	xargs := NewXargs()
	xargs.utilFactory = utilFactory
	xargs.utilArgs = args
	return xargs
}

// CLI invocation for *SomeXargs
func XargsCli(call []string) error {
	xargs := NewXargs()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := xargs.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return xargs.Exec(inPipe, outPipe, errPipe)
}
