package some

import (
	"bufio"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"time"
)

func init() {
	someutils.RegisterPipable(func() someutils.NamedPipable { return new(SomeTail) })
}

// SomeTail represents and performs a `tail` invocation
type SomeTail struct {
	Lines              int
	FollowByDescriptor bool
	FollowByName       bool
	SleepInterval      float64
	Filenames          []string
}

// Name() returns the name of the util
func (tail *SomeTail) Name() string {
	return "tail"
}

// ParseFlags parses flags from a commandline []string
func (tail *SomeTail) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("tail", "[options] [args...]", someutils.VERSION)
	flagSet.SetOutput(errPipe)

	flagSet.AliasedIntVar(&tail.Lines, []string{"n", "lines"}, 10, "number of lines to print")
	flagSet.AliasedFloat64Var(&tail.SleepInterval, []string{"s", "sleep"}, 1.0, "how long to sleep")
	//TODO!
	//flagSet.AliasedStringVar(&options.Follow, []string{"f", "follow"}, "", "follow (name|descriptor). Default is by descriptor (unsupported so far!!)")
	flagSet.BoolVar(&tail.FollowByName, "F", false, "follow by name")

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}
	tail.Filenames = flagSet.Args()
	return nil, 0
}

// Exec actually performs the tail
func (tail *SomeTail) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.AutoPipeErrInOut()
	invocation.AutoHandleSignals()
	if len(tail.Filenames) > 0 {
		for _, fileName := range tail.Filenames {
			finf, err := os.Stat(fileName)
			if err != nil {
				return err, 1
			}
			file, err := os.Open(fileName)
			if err != nil {
				return err, 1
			}
			seek := int64(0)
			if finf.Size() > 10000 {
				//just get last 10K (good enough for now)
				seek = finf.Size() - 10000
				_, err = file.Seek(seek, 0)
				if err != nil {
					return err, 1
				}
			}
			end, err := tailReader(file, seek, tail, invocation.OutPipe)
			if err != nil {
				file.Close()
				return err, 1
			}
			err = file.Close()
			if err != nil {
				return err, 1
			}
			if tail.FollowByName {
				sleepIntervalMs := time.Duration(tail.SleepInterval * 1000)
				for {
					//sleep n.x seconds
					//use milliseconds to get some accuracy with the int64
					time.Sleep(sleepIntervalMs * time.Millisecond)
					finf, err := os.Stat(fileName)
					if err != nil {
						return err, 1
					}
					file, err := os.Open(fileName)
					if err != nil {
						return err, 1
					}
					_, err = file.Seek(end, 0)
					if err != nil {
						return err, 1
					}
					if finf.Size() > end {
						end, err = tailReader(file, end, tail, invocation.OutPipe)
						if err != nil {
							file.Close()
							return err, 1
						}
					} else {
						//TODO start again
					}
					err = file.Close()
					if err != nil {
						return err, 1
					}
				}
			}
		}
	} else {
		//stdin ..
		_, err := tailReader(invocation.InPipe, 0, tail, invocation.OutPipe)
		if err != nil {
			return err, 1
		}
	}
	return nil, 0
}

func tailReader(file io.Reader, start int64, tail *SomeTail, out io.Writer) (int64, error) {
	var buffer []string
	end := start
	scanner := bufio.NewScanner(file)
	lastLine := tail.Lines - 1

	for scanner.Scan() {
		text := scanner.Text()
		end += int64(len(text) + 1) //for the \n character
		lastLine++
		if lastLine == tail.Lines {
			lastLine = 0
		}
		if lastLine >= len(buffer) {
			buffer = append(buffer, text)
		} else {
			buffer[lastLine] = text
		}
	}
	err := scanner.Err()
	if err != nil {
		return end, err
	}

	if lastLine == tail.Lines-1 {
		for _, text := range buffer {
			fmt.Fprintf(out, "%s\n", text)
		}
	} else {
		for _, text := range buffer[lastLine+1:] {
			fmt.Fprintf(out, "%s\n", text)
		}
		//if lastLine > 0 {
		for _, text := range buffer[:lastLine+1] {
			fmt.Fprintf(out, "%s\n", text)
		}
		//}
	}
	return end, nil
}

// Factory for *SomeTail
func NewTail() *SomeTail {
	return new(SomeTail)
}

// Factory for *SomeTail
func Tail(args ...string) *SomeTail {
	tail := NewTail()
	tail.Filenames = args
	return tail
}

// CLI invocation for *SomeTail
func TailCli(call []string) (error, int) {
	util := new(SomeTail)
	return someutils.StdInvoke((util), call)
}
