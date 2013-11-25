package someutils

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/laher/uggo"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

type LsOptions struct {
	LongList   bool
	Recursive  bool
	Human      bool
	AllFiles   bool
	OnePerLine bool
	Stdin      bool
}

var accessSymbols = "xwr"

func init() {
	Register(Util{
		"ls",
		Ls})
}

func Ls(call []string) error {
	options := LsOptions{}
	flagSet := uggo.NewFlagSetDefault("ls", "[options] [dirs...]", VERSION)
	flagSet.BoolVar(&options.LongList, "l", false, "Long, detailed listing")
	flagSet.AliasedBoolVar(&options.Recursive, []string{"R", "recursive"}, false, "Recurse into directories")
	flagSet.AliasedBoolVar(&options.Human, []string{"h", "human-readable"}, false, "Output sizes in a human readable format")
	flagSet.AliasedBoolVar(&options.AllFiles, []string{"a", "all"}, false, "Show all files (including dotfiles)")
	flagSet.BoolVar(&options.OnePerLine, "1", false, "One entry per line")
	flagSet.AliasedBoolVar(&options.Stdin, []string{"z", "stdin"}, false, "Read from stdin")

	out := tabwriter.NewWriter(os.Stdout, 4, 4, 1, ' ', 0)
	err := flagSet.Parse(call[1:])
	if err != nil {
		println("Error parsing flags")
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	args, err := getDirList(flagSet.Args(), options)
	if err != nil {
		return err
	}

	counter := 0
	lastWasDir := false
	for i, arg := range args {
		if !strings.HasPrefix(arg, ".") || options.AllFiles ||
			strings.HasPrefix(arg, "..") || "." == arg {
			argInfo, err := os.Stat(arg)
			if err != nil {
				fmt.Fprintln(os.Stderr, "stat failed for ", arg)
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
				if options.AllFiles {
					df, err := os.Stat(filepath.Dir(dir))
					if err != nil {
						fmt.Fprintf(out, "Error opening parent dir: %v", err)
					} else {
						printEntry("..", df, out, options, &counter)
					}
					df, err = os.Stat(dir)
					if err != nil {
						fmt.Fprintf(out, "Error opening dir: %v", err)
					} else {
						printEntry(".", df, out, options, &counter)
					}
				}

				err := list(out, dir, "", options, &counter)
				if err != nil {
					return err
				}
				if len(args) > 1 {
					fmt.Fprintf(out, "\n")
				}
			} else {

				listItem(argInfo, out, filepath.Dir(arg), "", options, &counter)
			}
			lastWasDir = argInfo.IsDir()
		}
	}
	out.Flush()
	return nil
}

func getDirList(globs []string, options LsOptions) ([]string, error) {
	if len(globs) <= 0 {
		if uggo.IsPipingStdin() {
			//check STDIN
			bio := bufio.NewReader(os.Stdin)
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

func list(out *tabwriter.Writer, dir, prefix string, options LsOptions, counter *int) error {
	if !strings.HasPrefix(dir, ".") || options.AllFiles ||
		strings.HasPrefix(dir, "..") || "." == dir {

		entries, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading dir '%s'", dir)
			return err
		}
		//dirs first, then files
		for _, entry := range entries {
			if entry.IsDir() {
				err = listItem(entry, out, dir, prefix, options, counter)
				if err != nil {
					return err
				}
			}
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				err = listItem(entry, out, dir, prefix, options, counter)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func listItem(entry os.FileInfo, out *tabwriter.Writer, dir, prefix string, options LsOptions, counter *int) error {
	if !strings.HasPrefix(entry.Name(), ".") || options.AllFiles {
		printEntry(entry.Name(), entry, out, options, counter)
		if entry.IsDir() && options.Recursive {
			folder := filepath.Join(prefix, entry.Name())
			if *counter%3 == 2 || options.LongList || options.OnePerLine {
				fmt.Fprintf(out, "%s:\n", folder)
			} else {
				fmt.Fprintf(out, "%s:\t", folder)
			}
			*counter += 1
			err := list(out, filepath.Join(dir, entry.Name()), folder, options, counter)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func printEntry(name string, e os.FileInfo, out *tabwriter.Writer, options LsOptions, counter *int) {
	if options.LongList {
		fmt.Fprintf(out, "%s\t", getModeString(e))
		if !e.IsDir() {
			fmt.Fprintf(out, "%s\t", getSizeString(e.Size(), options.Human))
		} else {
			fmt.Fprintf(out, "\t")
		}
		fmt.Fprintf(out, "%s\t", getModTimeString(e))
		//disabling due to native-only support
		//fmt.Fprintf(out, "%s\t", getUserString(e.Sys.(*syscall.Stat_t).Uid))
	}
	fmt.Fprintf(out, "%s%s\t", name, getEntryTypeString(e))
	if *counter%3 == 2 || options.LongList || options.OnePerLine {
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
