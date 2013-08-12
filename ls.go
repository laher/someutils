package someutils

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"text/tabwriter"
)

type LsOptions struct {
	longList   *bool
	recursive  *bool
	human      *bool
	allFiles   *bool
	onePerLine *bool
}

func Ls(call []string) error {
	options := LsOptions{}
	flagSet := flag.NewFlagSet("ls", flag.ContinueOnError)
	options.longList = flagSet.Bool("l", false, "Long, detailed listing")
	options.recursive = flagSet.Bool("r", false, "Recurse into directories")
	options.human = flagSet.Bool("h", false, "Output sizes in a human readable format")
	options.allFiles = flagSet.Bool("a", false, "Show all files (including dotfiles)")
	options.onePerLine = flagSet.Bool("1", false, "One entry per line")
	helpFlag := flagSet.Bool("help", false, "Show this help")
	out := tabwriter.NewWriter(os.Stdout, 4, 4, 1, ' ', 0)

	e := flagSet.Parse(splitSingleHyphenOpts(call[1:]))
	if e != nil {
		return e
	}

	if *helpFlag {
		println("`ls` [options] [dirs...]")
		flagSet.PrintDefaults()
		return nil
	}

	dirs, e := getDirList(flagSet.Args())
	if e != nil {
		return e
	}
	counter := 0
	for _, dir := range dirs {
		if *options.allFiles {
			//show . and ..
			df, err := os.Stat(filepath.Dir(dir))
			if err != nil {
				fmt.Fprintf(out, "Error opening parent dir", err)
			} else {
				printEntry("..", df, out, options, &counter)
			}
			df, err = os.Stat(dir)
			if err != nil {
				fmt.Fprintf(out, "Error opening dir", err)
			} else {
				printEntry(".", df, out, options, &counter)
			}
		}

		e := list(out, dir, "", options, &counter)
		if e != nil {
			return e
		}
	}
	out.Flush()
	return nil
}

func getDirList(args []string) ([]string, error) {
	if len(args) <= 0 {
		cwd, e := os.Getwd()
		return []string{cwd}, e
	}
	return args, nil
}

func list(out *tabwriter.Writer, dir, prefix string, options LsOptions, counter *int) error {
	entries, e := ioutil.ReadDir(dir)
	if e != nil {
		return e
	}
	for _, entry := range entries {
		printEntry(entry.Name(), entry, out, options, counter)
		if entry.IsDir() && *options.recursive {
			folder := filepath.Join(prefix, entry.Name())
			if *counter%3 == 2 || *options.longList || *options.onePerLine {
				fmt.Fprintf(out, "%s:\n", folder)
			} else {
				fmt.Fprintf(out, "%s:\t", folder)
			}
			*counter += 1
			e := list(out, filepath.Join(dir, entry.Name()), folder, options, counter)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func printEntry(name string, e os.FileInfo, out *tabwriter.Writer, options LsOptions, counter *int) {
	fmt.Fprintf(out, "%s%s\t", name, getEntryTypeString(e))
	if *options.longList {
		fmt.Fprintf(out, "%s\t", getModeString(e.Mode()))
		fmt.Fprintf(out, "%s\t", getSizeString(e.Size(), options.human))
		//fmt.Fprintf(out, "%s\t", getUserString(e.Sys.(*syscall.Stat_t).Uid))
	}
	if *counter%3 == 2 || *options.longList || *options.onePerLine {
		fmt.Fprintln(out, "")
	}
	*counter += 1
}

var accessSymbols = "xwr"

func getModeString(mode os.FileMode) (s string) {
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

func getSizeString(size int64, humanFlag *bool) (s string) {
	if !*humanFlag {
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
	return fmt.Sprintf("%7.3f%s", rSize, sizeSymbols[power:power+1])
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
