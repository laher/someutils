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
	someutils.RegisterPipable(func() someutils.CliPipable { return new(SomeWc) })
}

// SomeWc represents and performs a `wc` invocation
type SomeWc struct {
	IsBytes bool
	IsWords bool
	IsLines bool
	args    []string
}

// Name() returns the name of the util
func (wc *SomeWc) Name() string {
	return "wc"
}

// ParseFlags parses flags from a commandline []string
func (wc *SomeWc) ParseFlags(call []string, errWriter io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("wc", "[OPTION]... [FILE]...", someutils.VERSION)
	flagSet.SetOutput(errWriter)

	// TODO add flags here
	flagSet.AliasedBoolVar(&wc.IsLines, []string{"l", "lines"}, false, "Count lines")
	flagSet.AliasedBoolVar(&wc.IsWords, []string{"w", "words"}, false, "Count words")
	//	flagSet.AliasedBoolVar(&wc.IsChars, []string{"m", "chars"}, false, "Count characters")
	flagSet.AliasedBoolVar(&wc.IsBytes, []string{"c", "bytes"}, false, "Count bytes")

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}

	wc.args = flagSet.Args()
	return nil, 0
}

// Exec actually performs the wc
func (wc *SomeWc) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.AutoPipeErrInOut()
	invocation.AutoHandleSignals()
	if len(wc.args) > 0 {
		//treat no args as all args
		if !wc.IsWords && !wc.IsLines && !wc.IsBytes {
			wc.IsWords = true
			wc.IsLines = true
			wc.IsBytes = true
		}
		for _, fileName := range wc.args {
			bytes := int64(0)
			words := int64(0)
			lines := int64(0)
			//get byte count
			file, err := os.Open(fileName)
			if err != nil {
				return err, 1
			}
			err = countWords(file, wc, &bytes, &words, &lines)
			if err != nil {
				file.Close()
				return err, 1
			}
			err = file.Close()
			if err != nil {
				return err, 1
			}
			if wc.IsWords && !wc.IsLines && !wc.IsBytes {
				fmt.Fprintf(invocation.OutPipe, "%d %s\n", words, fileName)
			} else if !wc.IsWords && wc.IsLines && !wc.IsBytes {
				fmt.Fprintf(invocation.OutPipe, "%d %s\n", lines, fileName)
			} else if !wc.IsWords && !wc.IsLines && wc.IsBytes {
				fmt.Fprintf(invocation.OutPipe, "%d %s\n", bytes, fileName)
			} else {
				fmt.Fprintf(invocation.OutPipe, "%d %d %d %s\n", lines, words, bytes, fileName)
			}
		}
	} else {
		//stdin ..
		if !wc.IsWords && !wc.IsLines && !wc.IsBytes {
			wc.IsWords = true
		}
		bytes := int64(0)
		words := int64(0)
		lines := int64(0)
		err := countWords(invocation.InPipe, wc, &bytes, &words, &lines)
		if err != nil {
			return err, 1
		}
		if wc.IsWords && !wc.IsLines && !wc.IsBytes {
			fmt.Fprintf(invocation.OutPipe, "%d\n", words)
		} else if !wc.IsWords && wc.IsLines && !wc.IsBytes {
			fmt.Fprintf(invocation.OutPipe, "%d\n", lines)
		} else if !wc.IsWords && !wc.IsLines && wc.IsBytes {
			fmt.Fprintf(invocation.OutPipe, "%d\n", bytes)
		} else {
			fmt.Fprintf(invocation.OutPipe, "%d %d %d\n", lines, words, bytes)
		}
	}
	return nil, 0

}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func countWords(file io.Reader, wc *SomeWc, bytes *int64, words *int64, lines *int64) (err error) {
	lastWasSpace := false
	bio := bufio.NewReader(file)
	for err == nil {
		c, err := bio.ReadByte()
		if err != nil {
			if io.EOF == err {
				return nil
			}
			return err
		}
		*bytes += 1
		if isSpace(c) {
			if !lastWasSpace {
				*words += 1
			}
			lastWasSpace = true
		} else {
			lastWasSpace = false
		}
		if c == '\n' {
			*lines += 1
		}

	}
	return err
}

// Factory for *SomeWc
func NewWc() *SomeWc {
	return new(SomeWc)
}

// Factory for *SomeWc
func Wc(args ...string) *SomeWc {
	wc := NewWc()
	wc.args = args
	return wc
}

// CLI invocation for *SomeWc
func WcCli(call []string) (error, int) {
	util := new(SomeWc)
	return someutils.StdInvoke((util), call)
}
