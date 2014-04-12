package some

import (
	"bufio"
	"errors"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
)

func init() {
	someutils.RegisterPipable(func() someutils.CliPipable { return new(SomeXargs) })
}

// SomeXargs represents and performs a `xargs` invocation
type SomeXargs struct {
	utilFactory  someutils.CliPipableFactory
	utilArgs     []string
	maxProcesses int
}

// Name() returns the name of the util
func (xargs *SomeXargs) Name() string {
	return "xargs"
}

// ParseFlags parses flags from a commandline []string
func (xargs *SomeXargs) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("xargs", "[options] [args...]", someutils.VERSION)
	flagSet.SetOutput(errPipe)
	// TODO multiple processes at once ?
	flagSet.AliasedIntVar(&xargs.maxProcesses, []string{"P", "max-procs"}, 1, "Maximum processes")
	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}

	args := flagSet.Args()
	if len(args) < 1 {
		return errors.New("No command specified"), 1
	}
	if !someutils.Exists(args[0]) {
		return errors.New("Command does not exist."), 1
	}
	xargs.utilFactory = someutils.GetCliPipableFactory(args[0])
	xargs.utilArgs = args
	return nil, 0
}

// Exec actually performs the xargs
func (xargs *SomeXargs) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.AutoPipeErrInOut()
	invocation.AutoHandleSignals()
	util := xargs.utilFactory()
	args := xargs.newArgset(util.Name())
	reader := bufio.NewReader(invocation.InPipe)
	cont := true
	count := 0
	maxCount := 5
	for cont {
		if count >= maxCount {
			count = 0
			//fmt.Fprintf(errPipe, "args for '%s': %v\n", util.Name(), args)
			err, code := util.ParseFlags(args, invocation.ErrOutPipe)
			if err != nil {
				return err, code
			}
			err, code = util.Invoke(invocation)
			if err != nil {
				return err, code
			}
		}
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			cont = false
		} else if err != nil {
			return err, 1
		} else {
			args = append(args, string(line))
			if err != nil {
				return err, 1
			}
			count++
		}
	}
	//still more args to process
	if count > 0 {
		//fmt.Fprintf(errPipe, "args for '%s': %v\n", util.Name(), args)
		err, code := util.ParseFlags(args, invocation.ErrOutPipe)
		if err != nil {
			return err, code
		}
		err, code = util.Invoke(invocation)
		return err, code
	}
	return nil, 0
}

func (xargs *SomeXargs) newArgset(cmdName string) []string {
	args := []string{cmdName}
	args = append(args, xargs.utilArgs...)
	return args
}

// Factory for *SomeXargs
func NewXargs() *SomeXargs {
	return new(SomeXargs)
}

// Factory for *SomeXargs
func Xargs(utilFactory someutils.CliPipableFactory, args ...string) *SomeXargs {
	xargs := NewXargs()
	xargs.utilFactory = utilFactory
	xargs.utilArgs = args
	return xargs
}

// CLI invocation for *SomeXargs
func XargsCli(call []string) (error, int) {
	util := new(SomeXargs)
	return someutils.StdInvoke(util, call)
}
