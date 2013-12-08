package someutils

import (
	"fmt"
	"io"
	"bufio"
	"github.com/laher/uggo"
	"os"
)


type WcOptions struct {
	IsBytesOnly bool
	IsWordsOnly bool
	IsLinesOnly bool
}

func init() {
	Register(Util{
		"wc",
		Wc})
}

func Wc(call []string) error {
	options := WcOptions{}
	flagSet := uggo.NewFlagSetDefault("wc", "[options]", VERSION)
	err := flagSet.Parse(call[1:])
	flagSet.AliasedBoolVar(&options.IsLinesOnly, []string{"l", "lines"}, false, "Count lines")
	flagSet.AliasedBoolVar(&options.IsWordsOnly, []string{"w", "words"}, false, "Count words")
	flagSet.AliasedBoolVar(&options.IsBytesOnly, []string{"m", "chars"}, false, "Count characters")
	flagSet.AliasedBoolVar(&options.IsBytesOnly, []string{"c", "bytes"}, false, "Count bytes")

	if err != nil {
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	if len(flagSet.Args()) > 0 {
		for _, fileName := range flagSet.Args() {
			bytes := int64(-1)
			words := int64(-1)
			lines := int64(-1)
			//get byte count
			if !options.IsWordsOnly && !options.IsLinesOnly {
				finf, err := os.Stat(fileName)
				if err != nil {
					return err
				}
				bytes = finf.Size()
			}
			
			if !options.IsBytesOnly {
				file, err := os.Open(fileName)
				if err != nil {
					return err
				}
				err = wc(file, options, &bytes, &words, &lines)
				if err != nil {
					file.Close()
					return err
				}
				err = file.Close()
				if err != nil {
					return err
				}
				if options.IsWordsOnly {
					fmt.Println("%d %s", words, fileName)
				} else if options.IsLinesOnly {
					fmt.Println("%d %s", lines, fileName)
				} else {
					fmt.Println("%d %d %d %s", lines, words, bytes, fileName)
				}
			} else {
				fmt.Println("%d %s", bytes, fileName)
			}
		}
	} else {
		//stdin ..
		bytes := int64(-1)
		words := int64(-1)
		lines := int64(-1)
		return wc(os.Stdin, options, &bytes, &words, &lines)
	}
	if err != nil {
		return err
	}
	//println(wd)
	return nil
}

func wc(file io.Reader, options WcOptions, bytes *int64, words *int64, lines *int64) error {
	if !options.IsLinesOnly {
		//get lines ...
	}
	scanner := bufio.NewScanner(file)
	if options.IsLinesOnly {
		//ok
	} else {
		scanner.Split(bufio.ScanWords)
	}
	// Count the words.
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
		return err
	}
	return nil
}