package someutils

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	if !strings.Contains(link, ":") {
		link = "http://" + link
	}
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	filename := filepath.Base(resp.Request.URL.Path)
	//invalid filenames ...
	if filename == "" || filename == "/" ||filename == "\\" || filename == "." {
		filename = "index"
	}
	if !strings.Contains(filename, ".") {
		ct := resp.Header.Get("Content-Type")
		//println(ct)
		ext := "htm"
		mediatype, _, err := mime.ParseMediaType(ct)
		if err != nil {
			fmt.Printf("mime error: %v\n", err)
		} else {
			fmt.Printf("mime type: %v (from Content-Type %v)\n", mediatype, ct)
			slash := strings.Index(mediatype, "/")
			if slash != -1 {
				_, sub := mediatype[:slash], mediatype[slash+1:]
				if sub != "" {
					ext = sub
				}
			}
		}
		filename = filename + "." + ext
	}
	fmt.Printf("filename: %v\n", filename)
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	n, err := io.Copy(out, resp.Body)
	fmt.Printf("%v bytes written to %v\n",n, filename)
	if err != nil {
		return err
	}
	err = out.Close()
	return err
	//return errors.New("Not implemented")
}
