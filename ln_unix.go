//+build !windows

package someutils

import (
	"errors"
	"github.com/laher/uggo"
	"io"
	"os"
)

type SomeLn struct {
	target string
	linkName string
	IsForce    bool
	IsSymbolic bool
}
/*
func init() {
	Register({
		"ln",
		Ln})
}
*/
func (ln *SomeLn) ParseFlags(call []string, errPipe io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("ln", "[options] TARGET LINK_NAME", VERSION)
	flagSet.BoolVar(&ln.IsSymbolic, "s", false, "Symbolic")
	flagSet.BoolVar(&ln.IsForce, "f", false, "Force")

	e := flagSet.Parse(call[1:])
	if e != nil {
		return e
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	args := flagSet.Args()
	if len(args) < 2 {
		flagSet.Usage()
		return errors.New("Not enough args!")
	}
	ln.target = args[0]
	ln.linkName = args[1]
	return nil
}

func (ln *SomeLn) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	if ln.IsSymbolic {
		return os.Symlink(ln.target, ln.linkName)
	} else {
		return os.Link(ln.target, ln.linkName)
	}
}

