package someutils

import (
	"bufio"
	"fmt"
	"github.com/laher/uggo"
	"os"
)

type HeadOptions struct {
	lines     int
}

func init() {
	Register(Util{
		"head",
		Head})
}

func Head(call []string) error {
	options := HeadOptions{}
	flagSet := uggo.NewFlagSetDefault("head", "[options] [files...]", VERSION)
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
			if file, err := os.Open(fileName); err == nil {
				scanner := bufio.NewScanner(file)
				line := 1
				for scanner.Scan() && line <= options.lines {
					text := scanner.Text()
					fmt.Fprintf(os.Stdout, "%s\n", text)
					line++
				}
				err := scanner.Err()
				if err != nil {
					return err
				}
				file.Close()
			} else {
				return err
			}
		}
	} else {
		//stdin ..
		scanner := bufio.NewScanner(os.Stdin)
		line := 1
		for scanner.Scan() && line <= options.lines {
			text := scanner.Text()
			fmt.Fprintf(os.Stdout, "%s\n", text)
			line++
		}
		err := scanner.Err()
		if err != nil {
			return err
		}
	}
	return nil
}
