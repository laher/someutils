package someutils

import (
	"bufio"
	"fmt"
	"io"
	"github.com/laher/uggo"
	"os"
)

type TailOptions struct {
	lines int
}

func init() {
	Register(Util{
		"tail",
		Tail})
}

func Tail(call []string) error {
	options := TailOptions{}
	flagSet := uggo.NewFlagSetDefault("tail", "[options] [files...]", VERSION)
	flagSet.AliasedIntVar(&options.lines, []string{"n", "lines"}, 10, "number of lines to print")
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
			file, err := os.Open(fileName)
			if err != nil {
				return err
			}
			err = tail(file, options)
			if err != nil {
				file.Close()
				return err
			}
			err = file.Close()
			if err != nil {
				return err
			}
		}
	} else {
		//stdin ..
		return tail(os.Stdin, options)
	}
	return nil
}

func tail(file io.Reader, options TailOptions) error {
	var buffer []string
	scanner := bufio.NewScanner(file)
	lastLine := options.lines - 1
	for scanner.Scan() {
		text := scanner.Text()
		lastLine++
		if lastLine == options.lines {
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
		return err
	}

	//fmt.Fprintf(os.Stdout, "%s\n", text)
	if lastLine == options.lines-1 {
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
	return nil
}
