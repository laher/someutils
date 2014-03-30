package some

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
)

func init() {
	someutils.RegisterPipable(func() someutils.PipableCliUtil { return NewUnzip() })
}

// SomeUnzip represents and performs a `unzip` invocation
type SomeUnzip struct {
	// TODO: add members here
	destDir string
	isTest  bool

	zipname string
	files   []string
}

// Name() returns the name of the util
func (unzip *SomeUnzip) Name() string {
	return "unzip"
}

// TODO: add validation here

// ParseFlags parses flags from a commandline []string
func (unzip *SomeUnzip) ParseFlags(call []string, errWriter io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("unzip", "[options] file.zip [list...]", someutils.VERSION)
	flagSet.SetOutput(errWriter)
	destDir := "."
	flagSet.StringVar(&unzip.destDir, "d", destDir, "destination directory")
	test := false
	flagSet.BoolVar(&unzip.isTest, "t", test, "test archive data")

	// TODO add flags here

	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errWriter, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}

	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	args := flagSet.Args()
	if len(args) < 1 {
		return errors.New("No zip filename given")
	}
	unzip.zipname = args[0]
	unzip.files = args[1:]

	return nil
}

// Exec actually performs the unzip
func (unzip *SomeUnzip) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	if unzip.isTest {
		err := TestItems(unzip.zipname, unzip.files, outPipe, errPipe)
		if err != nil {
			return err
		}
	} else {
		err := UnzipItems(unzip.zipname, unzip.destDir, unzip.files, errPipe)
		if err != nil {
			return err
		}
	}
	return nil
}

func containsGlob(haystack []string, needle string, errPipe io.Writer) bool {
	for _, item := range haystack {
		m, err := filepath.Match(item, needle)
		if err != nil {
			fmt.Fprintf(errPipe, "Glob error %v", err)
			return false
		}
		if m == true {
			return true
		}
	}
	return false
}

func TestItems(zipfile string, includeFiles []string, outPipe io.Writer, errPipe io.Writer) error {
	r, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		flags := f.FileHeader.Flags
		if len(includeFiles) == 0 || containsGlob(includeFiles, f.Name, errPipe) {
			if flags&1 == 1 {
				fmt.Fprintf(outPipe, "[Password Protected:] %s\n", f.Name)
			} else {
				fmt.Fprintf(outPipe, "%s\n", f.Name)
			}
		}
	}
	return nil
}

func UnzipItems(zipfile, destDir string, includeFiles []string, errPipe io.Writer) error {

	r, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer r.Close()

	dinf, err := os.Stat(destDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} else {
			//doesnt exist
			err = os.MkdirAll(destDir, 0777) //TODO review permissions
			if err != nil {
				return err
			}
		}
	} else {
		if !dinf.IsDir() {
			return errors.New("destination is an existing non-directory")
		}
	}

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		finf := f.FileHeader.FileInfo()
		flags := f.FileHeader.Flags
		if flags&1 == 1 {
			fmt.Fprintf(errPipe, "WARN: Skipping password protected file (flags %v, '%s')\n", flags, f.Name)
		} else {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			destFileName := filepath.Join(destDir, f.Name)
			if finf.IsDir() {
				//mkdir ...
				fdinf, err := os.Stat(destFileName)
				if err != nil {
					if !os.IsNotExist(err) {
						return err
					} else {
						//doesnt exist
						err = os.MkdirAll(destFileName, finf.Mode())
						if err != nil {
							return err
						}
					}
				} else {
					if !fdinf.IsDir() {
						return errors.New("destination " + destFileName + " is an existing non-directory")
					}
				}
			} else {
				fileDestDir := filepath.Dir(destFileName)
				if fileDestDir != destDir {
					fdinf, err := os.Stat(fileDestDir)
					if err != nil {
						if !os.IsNotExist(err) {
							return err
						} else {
							//doesnt exist
							err = os.MkdirAll(fileDestDir, 0777) //TODO review dir permissions
							if err != nil {
								return err
							}
						}
					} else {
						if !fdinf.IsDir() {
							return errors.New("destination " + fileDestDir + " is an existing non-directory")
						}
					}
				}
				//TODO remove on error
				destFile, err := os.OpenFile(destFileName, os.O_CREATE, finf.Mode())
				defer destFile.Close()
				if err != nil {
					return err
				}
				_, err = io.Copy(destFile, rc)
				if err != nil {
					return err
				}

			}
			rc.Close()
		}
	}
	return nil
}

// Factory for *SomeUnzip
func NewUnzip() *SomeUnzip {
	return new(SomeUnzip)
}

// Factory for *SomeUnzip
func Unzip(zipname string, files ...string) *SomeUnzip {
	unzip := NewUnzip()
	unzip.zipname = zipname
	unzip.files = files
	return unzip
}

// CLI invocation for *SomeUnzip
func UnzipCli(call []string) error {
	unzip := NewUnzip()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := unzip.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return unzip.Exec(inPipe, outPipe, errPipe)
}
