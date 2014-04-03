package some

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

func init() {
	someutils.RegisterPipable(func() someutils.NamedPipable { return NewLs() })
}

// SomeLs represents and performs a `ls` invocation
type SomeLs struct {
	LongList   bool
	Recursive  bool
	Human      bool
	AllFiles   bool
	OnePerLine bool
	Stdin      bool

	globs []string
}

var accessSymbols = "xwr"

// Name() returns the name of the util
func (ls *SomeLs) Name() string {
	return "ls"
}

// TODO: add validation here

// ParseFlags parses flags from a commandline []string
func (ls *SomeLs) ParseFlags(call []string, errPipe io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("ls", "[options] [dirs...]", someutils.VERSION)
	flagSet.SetOutput(errPipe)

	flagSet.BoolVar(&ls.LongList, "l", false, "Long, detailed listing")
	flagSet.AliasedBoolVar(&ls.Recursive, []string{"R", "recursive"}, false, "Recurse into directories")
	flagSet.AliasedBoolVar(&ls.Human, []string{"h", "human-readable"}, false, "Output sizes in a human readable format")
	flagSet.AliasedBoolVar(&ls.AllFiles, []string{"a", "all"}, false, "Show all files (including dotfiles)")
	flagSet.BoolVar(&ls.OnePerLine, "1", false, "One entry per line")
	flagSet.AliasedBoolVar(&ls.Stdin, []string{"z", "stdin"}, false, "Read from stdin")

	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errPipe, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}

	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	//fmt.Fprintf(errPipe, "ls args: %+v\n", flagSet.Args())
	ls.globs = flagSet.Args()
	return nil
}

// Exec actually performs the ls
func (ls *SomeLs) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	out := tabwriter.NewWriter(outPipe, 4, 4, 1, ' ', 0)
	args, err := getDirList(ls.globs, ls, inPipe, outPipe, errPipe)
	if err != nil {
		return err
	}

	counter := 0
	lastWasDir := false
	for i, arg := range args {
		if !strings.HasPrefix(arg, ".") || ls.AllFiles ||
			strings.HasPrefix(arg, "..") || "." == arg {
			argInfo, err := os.Stat(arg)
			if err != nil {
				fmt.Fprintln(errPipe, "stat failed for ", arg)
				return err
			}
			if argInfo.IsDir() {
				if len(args) > 1 { //if more than one, print dir name before contents
					if i > 0 {
						fmt.Fprintf(out, "\n")
					}
					if !lastWasDir {
						fmt.Fprintf(out, "\n")
					}
					fmt.Fprintf(out, "%s:\n", arg)
				}
				dir := arg

				//show . and ..
				if ls.AllFiles {
					df, err := os.Stat(filepath.Dir(dir))
					if err != nil {
						fmt.Fprintf(out, "Error opening parent dir: %v", err)
					} else {
						printEntry("..", df, out, ls, &counter)
					}
					df, err = os.Stat(dir)
					if err != nil {
						fmt.Fprintf(out, "Error opening dir: %v", err)
					} else {
						printEntry(".", df, out, ls, &counter)
					}
				}

				err := list(out, errPipe, dir, "", ls, &counter)
				if err != nil {
					return err
				}
				if len(args) > 1 {
					fmt.Fprintf(out, "\n")
				}
			} else {

				listItem(argInfo, out, errPipe, filepath.Dir(arg), "", ls, &counter)
			}
			lastWasDir = argInfo.IsDir()
		}
	}
	out.Flush()
	return nil

}

// Factory for *SomeLs
func NewLs() *SomeLs {
	return new(SomeLs)
}

func LsFactory() someutils.PipableCliUtil {
	return NewLs()
}

// Factory for *SomeLs
func Ls(args ...string) *SomeLs {
	ls := NewLs()
	ls.globs = args
	return ls
}

// CLI invocation for *SomeLs
func LsCli(call []string) error {
	ls := NewLs()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := ls.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return ls.Exec(inPipe, outPipe, errPipe)
}

func getDirList(globs []string, ls *SomeLs, inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) ([]string, error) {
	if len(globs) <= 0 {
		if uggo.IsPipingStdin() {
			//check STDIN
			bio := bufio.NewReader(inPipe)
			//defer bio.Close()
			line, hasMoreInLine, err := bio.ReadLine()
			if err == nil {
				//adding from stdin
				globs = append(globs, strings.TrimSpace(string(line)))
			} else {
				//ok
			}
			for hasMoreInLine {
				line, hasMoreInLine, err = bio.ReadLine()
				if err == nil {
					//adding from stdin
					globs = append(globs, string(line))
				} else {
					//finish
				}
			}
		} else {
			//NOT piping. Just use cwd by default.
			cwd, err := os.Getwd()
			return []string{cwd}, err
		}
	}

	args := []string{}
	for _, glob := range globs {
		results, err := filepath.Glob(glob)
		if err != nil {
			return args, err
		}
		if len(results) < 1 { //no match
			return args, errors.New("ls: cannot access " + glob + ": No such file or directory")
		}
		args = append(args, results...)
	}
	return args, nil
}

func list(out *tabwriter.Writer, errPipe io.Writer, dir, prefix string, ls *SomeLs, counter *int) error {
	if !strings.HasPrefix(dir, ".") || ls.AllFiles ||
		strings.HasPrefix(dir, "..") || "." == dir {

		entries, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Fprintf(errPipe, "Error reading dir '%s'", dir)
			return err
		}
		//dirs first, then files
		for _, entry := range entries {
			if entry.IsDir() {
				err = listItem(entry, out, errPipe, dir, prefix, ls, counter)
				if err != nil {
					return err
				}
			}
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				err = listItem(entry, out, errPipe, dir, prefix, ls, counter)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func listItem(entry os.FileInfo, out *tabwriter.Writer, errPipe io.Writer, dir, prefix string, ls *SomeLs, counter *int) error {
	if !strings.HasPrefix(entry.Name(), ".") || ls.AllFiles {
		printEntry(entry.Name(), entry, out, ls, counter)
		if entry.IsDir() && ls.Recursive {
			folder := filepath.Join(prefix, entry.Name())
			if *counter%3 == 2 || ls.LongList || ls.OnePerLine {
				fmt.Fprintf(out, "%s:\n", folder)
			} else {
				fmt.Fprintf(out, "%s:\t", folder)
			}
			*counter += 1
			err := list(out, errPipe, filepath.Join(dir, entry.Name()), folder, ls, counter)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func printEntry(name string, e os.FileInfo, out *tabwriter.Writer, ls *SomeLs, counter *int) {
	if ls.LongList {
		fmt.Fprintf(out, "%s\t", getModeString(e))
		if !e.IsDir() {
			fmt.Fprintf(out, "%s\t", getSizeString(e.Size(), ls.Human))
		} else {
			fmt.Fprintf(out, "\t")
		}
		fmt.Fprintf(out, "%s\t", getModTimeString(e))
		//disabling due to native-only support
		//fmt.Fprintf(out, "%s\t", getUserString(e.Sys.(*syscall.Stat_t).Uid))
	}
	fmt.Fprintf(out, "%s%s\t", name, getEntryTypeString(e))
	if *counter%3 == 2 || ls.LongList || ls.OnePerLine {
		fmt.Fprintln(out, "")
	}
	*counter += 1
}

func getModTimeString(e os.FileInfo) (s string) {
	s = e.ModTime().Format("Jan 2 15:04")
	return
}

func getModeString(e os.FileInfo) (s string) {
	mode := e.Mode()
	if e.IsDir() {
		s = "d"
	} else {
		s = "-"
	}
	for i := 8; i >= 0; i-- {
		if mode&(1<<uint(i)) == 0 {
			s += "-"
		} else {
			char := i % 3
			s += accessSymbols[char : char+1]
		}
	}
	return
}

var sizeSymbols = "BkMGT"

func getSizeString(size int64, humanFlag bool) (s string) {
	if !humanFlag {
		return fmt.Sprintf("%9dB", size)
	}
	var power int
	if size == 0 {
		power = 0
	} else {
		power = int(math.Log(float64(size)) / math.Log(1024.0))
	}
	if power > len(sizeSymbols)-1 {
		power = len(sizeSymbols) - 1
	}
	rSize := float64(size) / math.Pow(1024, float64(power))
	return fmt.Sprintf("%7.1f%s", rSize, sizeSymbols[power:power+1])
}

func getEntryTypeString(e os.FileInfo) string {
	if e.IsDir() {
		return string(os.PathSeparator)
		/*	} else if e.IsBlock() {
				return "<>"
			} else if e.IsFifo() {
				return ">>"
			} else if e.IsSymlink() {
				return "@"
			} else if e.IsSocket() {
				return "&"
			} else if e.IsRegular() && (e.Mode&0001 == 0001) {
				return "*" */
	}
	return ""
}

func getUserString(id int) string {
	return fmt.Sprintf("%03d", id)
}
