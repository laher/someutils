package someutils

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
)

func init() {
	Register(Util{
		"unzip",
		Unzip})
}

func Unzip(call []string) error {

	flagSet := uggo.NewFlagSetDefault("unzip", "[options] file.zip [list]", VERSION)
	destDir := "."
	flagSet.StringVar(&destDir, "d", destDir, "destination directory")
	test := false
	flagSet.BoolVar(&test, "t", test, "test archive data")

	err := flagSet.Parse(call[1:])
	if err != nil {
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	args := flagSet.Args()
	if len(args) < 1 {
		return errors.New("No zip filename given")
	}
	zipname := args[0]
	files := args[1:]
	if test {
		err = TestItems(zipname, files)
		if err != nil {
			return err
		}
	} else {
		err = UnzipItems(zipname, destDir, files)
		if err != nil {
			return err
		}
	}
	return nil
}

func containsGlob(haystack []string, needle string) bool {
	for _, item := range haystack {
		m, err := filepath.Match(item, needle)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Glob error %v", err)
			return false
		}
		if m == true {
			return true
		}
	}
	return false
}

func TestItems(zipfile string, includeFiles []string) error {
	r, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		flags := f.FileHeader.Flags
		if len(includeFiles) == 0 || containsGlob(includeFiles, f.Name) {
			if flags&1 == 1 {
				fmt.Printf("[Password Protected:] %s\n", f.Name)
			} else {
				fmt.Printf("%s\n", f.Name)
			}
		}
	}
	return nil
}

func UnzipItems(zipfile, destDir string, includeFiles []string) error {

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
			fmt.Fprintf(os.Stderr, "WARN: Skipping password protected file (flags %v, '%s')\n", flags, f.Name)
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
