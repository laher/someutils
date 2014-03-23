package some

import (
	"archive/tar"
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
)

func init() {
	someutils.RegisterSome(func() someutils.SomeUtil { return NewTar() })
}

// SomeTar represents and performs a `tar` invocation
type SomeTar struct {
	IsCreate  bool
	IsList    bool
	IsExtract bool
	IsAppend  bool
	//IsCatenate bool
	IsVerbose       bool
	ArchiveFilename string

	args []string
}

// Name() returns the name of the util
func (tar *SomeTar) Name() string {
	return "tar"
}

// TODO: add validation here

// ParseFlags parses flags from a commandline []string
func (t *SomeTar) ParseFlags(call []string, errWriter io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("tar", "[option...] [FILE...]", someutils.VERSION)
	flagSet.SetOutput(errWriter)
	flagSet.AliasedBoolVar(&t.IsCreate, []string{"c", "create"}, false, "create a new archive")
	flagSet.AliasedBoolVar(&t.IsAppend, []string{"r", "append"}, false, "append files to the end of an archive")
	flagSet.AliasedBoolVar(&t.IsList, []string{"t", "list"}, false, "list the contents of an archive")
	flagSet.AliasedBoolVar(&t.IsExtract, []string{"x", "extract", "get"}, false, "extract files from an archive")
	flagSet.AliasedBoolVar(&t.IsVerbose, []string{"v", "verbose"}, false, "verbosely list files processed")
	flagSet.AliasedStringVar(&t.ArchiveFilename, []string{"f", "file"}, "", "use given archive file or device")

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
	if countTrue(t.IsCreate, t.IsAppend, t.IsList, t.IsExtract) != 1 {
		return errors.New("You must use *one* of -c, -t, -x, -r (create, list, extract or append), plus -f")
	}
	t.args = flagSet.Args()

	return nil
}

func countTrue(args ...bool) int {
	count := 0
	for _, arg := range args {
		if arg {
			count++
		}
	}
	return count
}

// Exec actually performs the tar
func (t *SomeTar) Exec(pipes someutils.Pipes) error {
	//overrideable??
	destDir := "."
	if t.IsCreate {
		//OK
		//fmt.Printf("Create %s\n", t.ArchiveFilename)
		err := TarItems(t.ArchiveFilename, t.args, t, pipes.Out())
		return err
	} else if t.IsAppend {
		//hmm is this OK with STDIN??
		if t.ArchiveFilename == "" {
			return errors.New("Filename (-f) must be provided in Append mode")
		}
		//OK
		//fmt.Printf("Append %s\n", t.ArchiveFilename)
		err := TarItems(t.ArchiveFilename, t.args, t, pipes.Out())
		return err
	} else if t.IsList {
		//fmt.Println("List", t.ArchiveFilename)
		err := TestTarItems(t.ArchiveFilename, t.args, pipes.In(), pipes.Out())
		return err
	} else if t.IsExtract {
		//fmt.Println("Extract", t.ArchiveFilename)
		err := UntarItems(t.ArchiveFilename, destDir, t.args, t, pipes.In(), pipes.Out())
		return err
	} else {
		return errors.New("You must use ONLY one of -c, -t or -x (create, list or extract), plus -f")
	}
}

func TarItems(tarfile string, includeFiles []string, t *SomeTar, outPipe io.Writer) error {
	var mode os.FileMode
	var zf *os.File
	var out io.Writer
	var err error
	if tarfile != "" {
		tinf, err := os.Stat(tarfile)
		mode = tinf.Mode()
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		}
		flags := os.O_RDWR
		if t.IsAppend {
			//flags = flags | os.O_APPEND
		} else {
			flags = flags | os.O_TRUNC | os.O_CREATE
		}
		zf, err = os.OpenFile(tarfile, flags, mode)
		if t.IsAppend {
			//println("append")
			//go to start of file footer
			zf.Seek(-1024, os.SEEK_END)
		}
		if err != nil {
			return err
		}
		defer zf.Close()
		out = zf
	} else {
		out = outPipe
	}

	zw := tar.NewWriter(out)
	defer zw.Close()

	for _, itemS := range includeFiles {
		//todo: relative/full path checking
		item := someutils.ArchiveItem{itemS, itemS, nil}
		err = addFileToTar(zw, item, t.IsVerbose, outPipe)
		if err != nil {
			return err
		}
	}
	//get error where possible
	err = zw.Close()
	if err != nil {
		return err
	}
	if tarfile != "" {
		err = zf.Close()
	}
	return err
}

func UntarItems(tarfile, destDir string, includeFiles []string, t *SomeTar, inPipe io.Reader, outPipe io.Writer) error {
	var in io.Reader
	if tarfile != "" {
		r, err := os.Open(tarfile)
		if err != nil {
			return err
		}
		defer r.Close()
		in = r
	} else {
		in = inPipe
	}
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

	tr := tar.NewReader(in)

	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return err
		}
		destFileName := filepath.Join(destDir, hdr.Name)
		finf := hdr.FileInfo()
		if finf.IsDir() {
			if t.IsVerbose {
				fmt.Fprintf(outPipe, "Making dir %s:\n", hdr.Name)
			}
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
			if t.IsVerbose {
				fmt.Fprintf(outPipe, "%s\n", hdr.Name)
			}
			//fmt.Printf("Contents of %s to %s\n", hdr.Name, destFileName)
			flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
			destFile, err := os.OpenFile(destFileName, flags, finf.Mode())
			defer destFile.Close()
			if err != nil {
				return err
			}
			_, err = io.Copy(destFile, tr)
			if err != nil {
				return err
			}
			err = destFile.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func TestTarItems(tarfile string, includeFiles []string, inPipe io.Reader, outPipe io.Writer) error {
	var in io.Reader
	if tarfile != "" {
		r, err := os.Open(tarfile)
		if err != nil {
			return err
		}
		defer r.Close()
		in = r
	} else {
		in = inPipe
	}
	tr := tar.NewReader(in)
	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return err
		}
		modeString := getModeString(hdr.FileInfo())
		modTimeString := getModTimeString(hdr.FileInfo())
		sizeString := getSizeString(hdr.FileInfo().Size(), false)
		fmt.Fprintf(outPipe, "%s %s %s %s\n", modeString, sizeString, modTimeString, hdr.Name)
	}
	return nil
}

func addFileToTar(zw *tar.Writer, item someutils.ArchiveItem, isVerbose bool, outPipe io.Writer) error {
	if isVerbose {
		fmt.Fprintf(outPipe, "Adding %s\n", item.FileSystemPath)
	}
	binfo, err := os.Stat(item.FileSystemPath)
	if err != nil {
		return err
	}
	if binfo.IsDir() {
		header, err := tar.FileInfoHeader(binfo, "")
		if err != nil {
			return err
		}
		header.Name = item.ArchivePath
		err = zw.WriteHeader(header)
		if err != nil {
			return err
		}
		file, err := os.Open(item.FileSystemPath)
		if err != nil {
			return err
		}
		fis, err := file.Readdir(0)
		for _, fi := range fis {
			err = addFileToTar(zw, someutils.ArchiveItem{filepath.Join(item.FileSystemPath, fi.Name()), filepath.Join(item.ArchivePath, fi.Name()), nil}, isVerbose, outPipe)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(binfo, "")
		if err != nil {
			return err
		}
		header.Name = item.ArchivePath
		err = zw.WriteHeader(header)
		if err != nil {
			return err
		}
		bf, err := os.Open(item.FileSystemPath)
		if err != nil {
			return err
		}
		defer bf.Close()
		_, err = io.Copy(zw, bf)
		if err != nil {
			return err
		}
		err = zw.Flush()
		if err != nil {
			return err
		}
		err = bf.Close()
		if err != nil {
			return err
		}
	}
	return err
}

// Factory for *SomeTar
func NewTar() *SomeTar {
	return new(SomeTar)
}

// Fluent factory for *SomeTar
func Tar(archiveFilename string, args ...string) *SomeTar {
	tar := NewTar()
	tar.ArchiveFilename = archiveFilename
	tar.args = args
	return tar
}

// CLI invocation for *SomeTar
func TarCli(call []string) error {
	tar := NewTar()
	pipes := someutils.StdPipes()
	err := tar.ParseFlags(call, pipes.Err())
	if err != nil {
		return err
	}
	return tar.Exec(pipes)
}
