package someutils

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"os"
	"bufio"
	"strings"
)

func init() {
	Register(Util{
		"grep",
		Grep})
}

type GrepOptions struct {
	IsPerl *bool
	IsExtended *bool
	IsIgnoreCase *bool
	IsInvertMatch *bool
	IsPrintFilename *bool
	IsPrintLineNumber *bool
	IsRecurse *bool
	IsQuiet *bool
	LinesBefore *int
	LinesAfter *int
	LinesAround *int
}

func Grep(call []string) error {

	options := GrepOptions{}
	flagSet := flag.NewFlagSet("grep", flag.ContinueOnError)
	options.IsPerl = flagSet.Bool("P", false, "Perl-style regex")
	options.IsExtended = flagSet.Bool("E", true, "Extended regex (default)")
	options.IsIgnoreCase = flagSet.Bool("i", false, "ignore case")
	options.IsPrintFilename = flagSet.Bool("H", true, "print the file name for each match")
	options.IsPrintLineNumber = flagSet.Bool("n", false, "print the line number for each match")
	options.IsInvertMatch = flagSet.Bool("v", false, "invert match")
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
	pattern := args[0]
	reg, err := compile(pattern, options)
	if err != nil {
		return err
	}

	globs := []string{}
	if len(args) > 1 {
		globs = args[1:]
		files := []string{}
		for _, glob := range globs {
			results, err := filepath.Glob(glob)
			if err != nil {
				return err
			}
			if len(results)<1 { //no match
				return errors.New("grep: cannot access "+glob+": No such file or directory")
			}
			files = append(files, results...)
		}
		return grep(reg, files, options)
	} else {
		if IsPipingStdin() {
			//check STDIN
			return grepReader(os.Stdin, "", reg, options)
			
		} else {
			//NOT piping.
			return errors.New("Not enough args")
		}
	}
}

func grep(reg *regexp.Regexp, files []string, options GrepOptions) error {
	
	for _,filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		err = grepReader(file, filename, reg, options)
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func grepReader(file io.Reader, filename string, reg *regexp.Regexp, options GrepOptions) error {
	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return err
		}
		line := scanner.Text()
		candidate := line
		if *options.IsIgnoreCase && !*options.IsPerl {
			candidate = strings.ToLower(line)
		}
		isMatch := reg.MatchString(candidate) 
		if (isMatch && !*options.IsInvertMatch) || (!isMatch && *options.IsInvertMatch) {
			if *options.IsPrintFilename && filename != "" {
				fmt.Printf("%s:", filename)
			}
			if *options.IsPrintLineNumber {
				fmt.Printf("%d:", lineNumber)
			}
			fmt.Println(line)
		}
		lineNumber+=1
	}
	return nil
}

func compile(pattern string, options GrepOptions) (*regexp.Regexp, error) {
	if *options.IsPerl {
		if *options.IsIgnoreCase && !strings.HasPrefix(pattern, "(?") {
			pattern = "(?i)"+pattern
		}
		return regexp.Compile(pattern)
	} else {
		if *options.IsIgnoreCase {
			pattern = strings.ToLower(pattern)
		}
		return regexp.CompilePOSIX(pattern)
	}
}