package someutils

import (
	"bufio"
	"fmt"
	"github.com/laher/uggo"
	"io"
	"os"
	"strings"
)

type CatOptions struct {
	showEnds     bool
	numbers      bool
	squeezeBlank bool
}

func init() {
	Register(Util{
		"cat",
		Cat})
}

func (o CatOptions) isStraightCopy() bool {
	return !o.showEnds && !o.numbers && !o.squeezeBlank
}

func Cat(call []string) error {
	options := CatOptions{}
	flagSet := uggo.NewFlagSetDefault("cat", "[options] [files...]", VERSION)
	flagSet.AliasedBoolVar(&options.showEnds, []string{"E", "show-ends"}, false, "display $ at end of each line")
	flagSet.AliasedBoolVar(&options.numbers, []string{"n", "number"}, false, "number all output lines")
	flagSet.AliasedBoolVar(&options.squeezeBlank, []string{"s", "squeeze-blank"}, false, "squeeze repeated empty output lines")

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
				//if straightCopy
				if options.isStraightCopy() {
					_, err = io.Copy(os.Stdout, file)
					if err != nil {
						return err
					}
				} else {
					scanner := bufio.NewScanner(file)
					line := 1
					var prefix string
					var suffix string
					for scanner.Scan() {
						text := scanner.Text()
						if !options.squeezeBlank || len(strings.TrimSpace(text)) > 0 {
							if options.numbers {
								prefix = fmt.Sprintf("%d ", line)
							} else {
								prefix = ""
							}
							if options.showEnds {
								suffix = "$"
							} else {
								suffix = ""
							}
							fmt.Fprintf(os.Stdout, "%s%s%s\n", prefix, text, suffix)

						}
						line++
					}
					err := scanner.Err()
					if err != nil {
						return err
					}
				}
				file.Close()
			} else {
				return err
			}
		}
	} else {
		_, err = io.Copy(os.Stdout, os.Stdin)
		if err != nil {
			return err
		}
	}
	return nil
}
