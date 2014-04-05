package some

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func init() {
	someutils.RegisterPipable(func() someutils.NamedPipable { return NewGrep() })
}

// SomeGrep represents and performs a `grep` invocation
type SomeGrep struct {
	IsPerl            bool
	IsExtended        bool
	IsIgnoreCase      bool
	IsInvertMatch     bool
	IsPrintFilename   bool
	IsPrintLineNumber bool
	IsRecurse         bool
	IsQuiet           bool
	LinesBefore       int
	LinesAfter        int
	LinesAround       int

	pattern string
	globs   []string
}

// Name() returns the name of the util
func (grep *SomeGrep) Name() string {
	return "grep"
}

// ParseFlags parses flags from a commandline []string
func (grep *SomeGrep) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("grep", "[options] PATTERN [files...]", someutils.VERSION)
	flagSet.SetOutput(errPipe)
	flagSet.AliasedBoolVar(&grep.IsPerl, []string{"P", "perl-regexp"}, false, "Perl-style regex")
	flagSet.AliasedBoolVar(&grep.IsExtended, []string{"E", "extended-regexp"}, true, "Extended regex (default)")
	flagSet.AliasedBoolVar(&grep.IsIgnoreCase, []string{"i", "ignore-case"}, false, "ignore case")
	flagSet.AliasedBoolVar(&grep.IsPrintFilename, []string{"H", "with-filename"}, true, "print the file name for each match")
	flagSet.AliasedBoolVar(&grep.IsPrintLineNumber, []string{"n", "line-number"}, false, "print the line number for each match")
	flagSet.AliasedBoolVar(&grep.IsInvertMatch, []string{"v", "invert-match"}, false, "invert match")
	// disable for now
	//	flagSet.AliasedBoolVar(&grep.IsRecurse, []string{"r", "recurse"}, false, "recurse into subdirectories")

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}
	args := flagSet.Args()
	if len(args) < 1 {
		flagSet.Usage()
		return errors.New("Not enough args"), 1
	}
	grep.pattern = args[0]

	if len(args) > 1 {
		grep.globs = args[1:]
	} else {
		grep.globs = []string{}
	}

	return nil, 0
}

// Exec actually performs the grep
func (grep *SomeGrep) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) (error, int) {
	reg, err := compile(grep.pattern, grep)
	if err != nil {
		return err, 1
	}
	if len(grep.globs) > 0 {
		files := []string{}
		for _, glob := range grep.globs {
			results, err := filepath.Glob(glob)
			if err != nil {
				return err, 1
			}
			if len(results) < 1 { //no match
				return errors.New("grep: cannot access " + glob + ": No such file or directory"), 1
			}
			files = append(files, results...)
		}
		err = grepAll(reg, files, grep, outPipe)
		if err != nil {
			return err, 1
		}
	} else {
		if uggo.IsPipingStdin() {
			//check STDIN
			err = grepReader(inPipe, "", reg, grep, outPipe)
			if err != nil {
				return err, 1
			}
		} else {
			//NOT piping.
			return errors.New("Not enough args"), 1
		}
	}
	return nil, 0
}

func grepAll(reg *regexp.Regexp, files []string, grep *SomeGrep, out io.Writer) error {
	for _, filename := range files {
		fi, err := os.Stat(filename)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			//recurse here
			if grep.IsRecurse {
				//
				fmt.Fprintf(out, "Recursion not implemented yet\n")
			}
		}
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		err = grepReader(file, filename, reg, grep, out)
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

func grepReader(file io.Reader, filename string, reg *regexp.Regexp, grep *SomeGrep, out io.Writer) error {
	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return err
		}
		line := scanner.Text()
		candidate := line
		if grep.IsIgnoreCase && !grep.IsPerl {
			candidate = strings.ToLower(line)
		}
		isMatch := reg.MatchString(candidate)
		if (isMatch && !grep.IsInvertMatch) || (!isMatch && grep.IsInvertMatch) {
			if grep.IsPrintFilename && filename != "" {
				fmt.Fprintf(out, "%s:", filename)
			}
			if grep.IsPrintLineNumber {
				fmt.Fprintf(out, "%d:", lineNumber)
			}
			fmt.Fprintln(out, line)
		}
		lineNumber += 1
	}
	return nil
}

func compile(pattern string, grep *SomeGrep) (*regexp.Regexp, error) {
	if grep.IsPerl {
		if grep.IsIgnoreCase && !strings.HasPrefix(pattern, "(?") {
			pattern = "(?i)" + pattern
		}
		return regexp.Compile(pattern)
	} else {
		if grep.IsIgnoreCase {
			pattern = strings.ToLower(pattern)
		}
		return regexp.CompilePOSIX(pattern)
	}
}

// Factory for *SomeGrep
func NewGrep() *SomeGrep {
	return new(SomeGrep)
}

// Factory for *SomeGrep
func Grep(args ...string) *SomeGrep {
	grep := NewGrep()
	grep.pattern = args[0]
	if len(args) > 1 {
		grep.globs = args[1:]
	} else {
		grep.globs = []string{}
	}
	return grep
}

// CLI invocation for *SomeGrep
func GrepCli(call []string) (error, int) {
	grep := NewGrep()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err, code := grep.ParseFlags(call, errPipe)
	if err != nil {
		return err, code
	}
	return grep.Exec(inPipe, outPipe, errPipe)
}
