package someutils

import (
	"bufio"
	"fmt"
	"github.com/laher/uggo"
	"io"
	"os"
)

type WcOptions struct {
	IsBytes bool
	IsWords bool
	IsLines bool
}

func init() {
	Register(Util{
		"wc",
		Wc})
}

func Wc(call []string) error {
	options := WcOptions{}
	flagSet := uggo.NewFlagSetDefault("wc", "[options]", VERSION)
	flagSet.AliasedBoolVar(&options.IsLines, []string{"l", "lines"}, false, "Count lines")
	flagSet.AliasedBoolVar(&options.IsWords, []string{"w", "words"}, false, "Count words")
	//	flagSet.AliasedBoolVar(&options.IsChars, []string{"m", "chars"}, false, "Count characters")
	flagSet.AliasedBoolVar(&options.IsBytes, []string{"c", "bytes"}, false, "Count bytes")
	err := flagSet.Parse(call[1:])

	if err != nil {
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	if len(flagSet.Args()) > 0 {
		//treat no args as all args
		if !options.IsWords && !options.IsLines && !options.IsBytes {
			options.IsWords = true
			options.IsLines = true
			options.IsBytes = true
		}
		for _, fileName := range flagSet.Args() {
			bytes := int64(0)
			words := int64(0)
			lines := int64(0)
			//get byte count
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
			if options.IsWords && !options.IsLines && !options.IsBytes {
				fmt.Printf("%d %s\n", words, fileName)
			} else if !options.IsWords && options.IsLines && !options.IsBytes {
				fmt.Printf("%d %s\n", lines, fileName)
			} else if !options.IsWords && !options.IsLines && options.IsBytes {
				fmt.Printf("%d %s\n", bytes, fileName)
			} else {
				fmt.Printf("%d %d %d %s\n", lines, words, bytes, fileName)
			}
		}
	} else {
		//stdin ..
		if !options.IsWords && !options.IsLines && !options.IsBytes {
			options.IsWords = true
		}
		bytes := int64(0)
		words := int64(0)
		lines := int64(0)
		err = wc(os.Stdin, options, &bytes, &words, &lines)
		if err != nil {
			return err
		}
		if options.IsWords && !options.IsLines && !options.IsBytes {
			fmt.Printf("%d\n", words)
		} else if !options.IsWords && options.IsLines && !options.IsBytes {
			fmt.Printf("%d\n", lines)
		} else if !options.IsWords && !options.IsLines && options.IsBytes {
			fmt.Printf("%d\n", bytes)
		} else {
			fmt.Printf("%d %d %d\n", lines, words, bytes)
		}
	}
	if err != nil {
		return err
	}
	//println(wd)
	return nil
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func wc(file io.Reader, options WcOptions, bytes *int64, words *int64, lines *int64) (err error) {
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
