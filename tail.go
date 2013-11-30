package someutils

import (
	"bufio"
	"fmt"
	"io"
	"github.com/laher/uggo"
	"os"
	"time"
)

type TailOptions struct {
	Lines int
	FollowByDescriptor bool
	FollowByName bool
	SleepInterval float64
}

func init() {
	Register(Util{
		"tail",
		Tail})
}

func Tail(call []string) error {
	options := TailOptions{}
	flagSet := uggo.NewFlagSetDefault("tail", "[options] [files...]", VERSION)
	flagSet.AliasedIntVar(&options.Lines, []string{"n", "lines"}, 10, "number of lines to print")
	flagSet.AliasedFloat64Var(&options.SleepInterval, []string{"s", "sleep"}, 1.0, "how long to sleep")
	//TODO!
	//flagSet.AliasedStringVar(&options.Follow, []string{"f", "follow"}, "", "follow (name|descriptor). Default is by descriptor (unsupported so far!!)")
	flagSet.BoolVar(&options.FollowByName, "F", false, "follow by name")
	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	if len(flagSet.Args()) > 0 {
		for _, fileName := range flagSet.Args() {
			finf, err := os.Stat(fileName)
			if err != nil {
				return err
			}
			file, err := os.Open(fileName)
			if err != nil {
				return err
			}
			seek := int64(0)
			if finf.Size() > 10000 {
				//just get last 10K (good enough for now)
				seek = finf.Size() - 10000
				_, err = file.Seek(seek, 0)
				if err != nil {
					return err
				}
			}
			end, err := tail(file, seek, options)
			if err != nil {
				file.Close()
				return err
			}
			err = file.Close()
			if err != nil {
				return err
			}
			if options.FollowByName {
				sleepIntervalMs := time.Duration(options.SleepInterval * 1000)
				for {
					//sleep n.x seconds
					//use milliseconds to get some accuracy with the int64
					time.Sleep(sleepIntervalMs * time.Millisecond)
					finf, err := os.Stat(fileName)
					if err != nil {
						return err
					}
					file, err := os.Open(fileName)
					if err != nil {
						return err
					}
					_, err = file.Seek(end, 0)
					if err != nil {
						return err
					}
					if finf.Size() > end {
						end, err = tail(file, end, options)
						if err != nil {
							file.Close()
							return err
						}
					} else {
						//TODO start again
					}
					err = file.Close()
					if err != nil {
						return err
					}
				}
			}
		}
	} else {
		//stdin ..
		_, err = tail(os.Stdin, 0, options)
		return err
	}
	return nil
}

func tail(file io.Reader, start int64, options TailOptions) (int64, error) {
	var buffer []string
	end := start
	scanner := bufio.NewScanner(file)
	lastLine := options.Lines - 1
	
	for scanner.Scan() {
		text := scanner.Text()
		end += int64(len(text) + 1) //for the \n character
		lastLine++
		if lastLine == options.Lines {
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

	//fmt.Fprintf(os.Stdout, "%s\n", text)
	if lastLine == options.Lines-1 {
		for _, r := range buffer {
			println(r)
		}
	} else {
		for _, r := range buffer[lastLine+1:] {
			println(r)
		}
		//if lastLine > 0 {
		for _, r := range buffer[:lastLine+1] {
			println(r)
		}
		//}
	}
	return end, nil
}
