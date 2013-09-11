package someutils

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	Register(Util{
		"wget",
		Wget})
}

type WgetOptions struct {
}

func Wget(call []string) error {

	//options := WgetOptions{}
	flagSet := flag.NewFlagSet("wget", flag.ContinueOnError)
/*
	options.IsPerl = flagSet.Bool("P", false, "Perl-style regex")
	options.IsExtended = flagSet.Bool("E", true, "Extended regex (default)")
	options.IsIgnoreCase = flagSet.Bool("i", false, "ignore case")
	options.IsPrintFilename = flagSet.Bool("H", true, "print the file name for each match")
	options.IsPrintLineNumber = flagSet.Bool("n", false, "print the line number for each match")
	options.IsInvertMatch = flagSet.Bool("v", false, "invert match")
*/
	helpFlag := flagSet.Bool("help", false, "Show this help")

	err := flagSet.Parse(splitSingleHyphenOpts(call[1:]))
	if err != nil {
		return err
	}
	if *helpFlag {
		println("`grep` [options] PATTERN [files...]")
		flagSet.PrintDefaults()
		return nil
	}
	args := flagSet.Args()
	if len(args) < 1 {
		flagSet.PrintDefaults()
		return errors.New("Not enough args")
	}
	if len(args) > 0 {
		links := args
		return wget(links)
	} else {
		if IsPipingStdin() {
			//check STDIN
			return wget([]string{})
		} else {
			//NOT piping.
			return errors.New("Not enough args")
		}
	}
}

func wget(links []string) error {
	for _, link := range links {
		err := wgetOne(link)
		if err != nil {
			return err
		}
	}
	return nil
}
func wgetOne(link string) error {
	resp, err := http.Get(link)
	defer resp.Body.Close()
	filename := filepath.Base(resp.Request.URL.Path)
	if filename == "" || filename == "/" || filename == "." {
		filename = "index.htm"
	}
	fmt.Printf("filename: %v\n", filename)
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	n, err := io.Copy(out, resp.Body)
	fmt.Printf("%v b written to %v\n",n, filename)
	if err != nil {
		return err
	}
	err = out.Close()
	return err
	//return errors.New("Not implemented")
}
