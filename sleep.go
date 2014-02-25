package someutils

import (
	"errors"
	"github.com/laher/uggo"
	"strconv"
	"time"
)

func init() {
	Register(Util{
		"sleep",
		Sleep})
}
//very basic for the moment. No removal of suffix
func Sleep(call []string) error {

	flagSet := uggo.NewFlagSetDefault("sleep", "", VERSION)

	err := flagSet.Parse(call[1:])
	if err != nil {
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	if len(flagSet.Args()) < 1 {
		return errors.New("sleep: missing operand")
	}
	arg := flagSet.Args()[0]
	last := arg[len(arg)-1:]
	_, err = strconv.Atoi(last)
	if err==nil {
		arg = arg + "s"
	}
	num := arg[:len(arg)-1]
	unit := arg[len(arg)-1:]

	sleepAmount, err := strconv.Atoi(num)
	if err != nil {
		return err
	}
	var unitDur time.Duration
	switch unit {
	case "d":
		unitDur = time.Hour * 24
	case "s":
		unitDur = time.Second
	case "m":
		unitDur = time.Minute
	case "h":
		unitDur = time.Hour
	default:
		return errors.New("Invalid time interval "+arg)
	}
	time.Sleep(time.Duration(sleepAmount) * unitDur)
	return nil
}
