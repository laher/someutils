package some

import (
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"strconv"
	"time"
)

func init() {
	someutils.RegisterSome(func() someutils.SomeUtil { return NewSleep() })
}

// SomeSleep represents and performs a `sleep` invocation
type SomeSleep struct {
	unit string
	amount int
}

// Name() returns the name of the util
func (sleep *SomeSleep) Name() string {
	return "sleep"
}


// ParseFlags parses flags from a commandline []string
func (sleep *SomeSleep) ParseFlags(call []string, errWriter io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("sleep", "", someutils.VERSION)
	flagSet.SetOutput(errWriter)

	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errWriter, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}

	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	arg := flagSet.Args()[0]
	last := arg[len(arg)-1:]
	_, err = strconv.Atoi(last)
	if err == nil {
		arg = arg + "s"
	}
	num := arg[:len(arg)-1]
	sleep.unit = arg[len(arg)-1:]

	sleep.amount, err = strconv.Atoi(num)
	if err != nil {
		return err
	}
	return nil
}

// Exec actually performs the sleep
func (sleep *SomeSleep) Exec(pipes someutils.Pipes) error {
	var unitDur time.Duration
	switch sleep.unit {
	case "d":
		unitDur = time.Hour * 24
	case "s":
		unitDur = time.Second
	case "m":
		unitDur = time.Minute
	case "h":
		unitDur = time.Hour
	default:
		return errors.New("Invalid time interval " + sleep.unit)
	}
	time.Sleep(time.Duration(sleep.amount) * unitDur)
	return nil

}

// Factory for *SomeSleep
func NewSleep() *SomeSleep {
	return new(SomeSleep)
}

// Fluent factory for *SomeSleep
func Sleep(amount int, unit string) *SomeSleep {
	sleep := NewSleep()
	sleep.unit = unit
	sleep.amount = amount
	return sleep
}

// CLI invocation for *SomeSleep
func SleepCli(call []string) error {
	sleep := NewSleep()
	pipes := someutils.StdPipes()
	err := sleep.ParseFlags(call, pipes.Err())
	if err != nil {
		return err
	}
	return sleep.Exec(pipes)
}
