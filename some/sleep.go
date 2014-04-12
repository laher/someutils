package some

import (
	"errors"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"strconv"
	"time"
)

func init() {
	someutils.RegisterPipable(func() someutils.NamedPipable { return new(SomeSleep) })
}

// SomeSleep represents and performs a `sleep` invocation
type SomeSleep struct {
	unit   string
	amount int
}

// Name() returns the name of the util
func (sleep *SomeSleep) Name() string {
	return "sleep"
}

// ParseFlags parses flags from a commandline []string
func (sleep *SomeSleep) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("sleep", "", someutils.VERSION)
	flagSet.SetOutput(errPipe)

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
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
		return err, 1
	}
	return nil, 0
}

// Exec actually performs the sleep
func (sleep *SomeSleep) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.AutoPipeErrInOut()
	invocation.AutoHandleSignals()
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
		return errors.New("Invalid time interval " + sleep.unit), 1
	}
	time.Sleep(time.Duration(sleep.amount) * unitDur)
	return nil, 0

}

// Factory for *SomeSleep
func NewSleep() *SomeSleep {
	return new(SomeSleep)
}

// Factory for *SomeSleep
func Sleep(amount int, unit string) *SomeSleep {
	sleep := NewSleep()
	sleep.unit = unit
	sleep.amount = amount
	return sleep
}

// CLI invocation for *SomeSleep
func SleepCli(call []string) (error, int) {

	util := new(SomeSleep)
	return someutils.StdInvoke((util), call)

}
