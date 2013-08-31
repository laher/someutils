package someutils

import (
	"flag"
	"io"
	"os"
)

func init() {
	Register(Util{
		"cat",
		Cat})
}

type CatOptions struct {
	showEnds *bool
	numbers *bool
	squeezeBlank *bool
}

func Cat(call []string) error {

	options := CatOptions{}
	flagSet := flag.NewFlagSet("cat", flag.ContinueOnError)
	options.showEnds = flagSet.Bool("E", false, "display $ at end of each line")
	options.numbers = flagSet.Bool("n", false, "number all output lines")
	options.squeezeBlank = flagSet.Bool("s", false, "squeeze repeated empty output lines")
	helpFlag := flagSet.Bool("help", false, "Show this help")

	err := flagSet.Parse(splitSingleHyphenOpts(call[1:]))
	if err != nil {
		return err
	}

	if *helpFlag {
		println("`cat` [options] [files...]")
		flagSet.PrintDefaults()
		return nil
	}
	
	if len(flagSet.Args()) > 0 {
		for _, fileName := range flagSet.Args() {
			if file, err := os.Open(fileName); err == nil {
				_, err = io.Copy(os.Stdout, file)
				if err != nil {
					return err
				}
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