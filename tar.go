package someutils

import (
	"archive/tar"
	"errors"
	"fmt"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
)

type TarOptions struct {
	IsCreate  bool
	IsList    bool
	IsExtract bool
	IsAppend  bool
	//IsCatenate bool
	IsVerbose       bool
	ArchiveFilename string
}

func init() {
	Register(Util{
		"tar",
		Tar})
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

func Tar(call []string) error {
	options := TarOptions{}
	flagSet := uggo.NewFlagSetDefault("tar", "[OPTION...] [FILE]...", VERSION)
	flagSet.AliasedBoolVar(&options.IsCreate, []string{"c", "create"}, false, "create a new archive")
	flagSet.AliasedBoolVar(&options.IsAppend, []string{"r", "append"}, false, "append files to the end of an archive")
	flagSet.AliasedBoolVar(&options.IsList, []string{"t", "list"}, false, "list the contents of an archive")
	flagSet.AliasedBoolVar(&options.IsExtract, []string{"x", "extract", "get"}, false, "extract files from an archive")
	flagSet.AliasedBoolVar(&options.IsVerbose, []string{"v", "verbose"}, false, "verbosely list files processed")
	flagSet.AliasedStringVar(&options.ArchiveFilename, []string{"f", "file"}, "", "use given archive file or device")
	destDir := "."

	err := flagSet.Parse(call[1:])
	if err != nil {
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	//if !options.IsCreate && !options.IsList && !options.IsExtract {
	if countTrue(options.IsCreate, options.IsAppend, options.IsList, options.IsExtract) != 1 {
		return errors.New("You must use *one* of -c, -t, -x, -r (create, list, extract or append), plus -f")
	}
	args := flagSet.Args()

	if options.IsCreate {
		//OK
		//fmt.Printf("Create %s\n", options.ArchiveFilename)
		err = TarItems(options.ArchiveFilename, args, options)
		if err != nil {
			return err
		}
	} else if options.IsAppend {
		//hmm is this OK with STDIN??
		if options.ArchiveFilename == "" {
			return errors.New("Filename (-f) must be provided in Append mode")
		}
		//OK
		//fmt.Printf("Append %s\n", options.ArchiveFilename)
		err = TarItems(options.ArchiveFilename, args, options)
		if err != nil {
			return err
		}
	} else if options.IsList {
		//fmt.Println("List", options.ArchiveFilename)
		err = TestTarItems(options.ArchiveFilename, args)
		if err != nil {
			return err
		}

	} else if options.IsExtract {
		//fmt.Println("Extract", options.ArchiveFilename)
		err = UntarItems(options.ArchiveFilename, destDir, args, options)
		if err != nil {
			return err
		}
	} else {
		return errors.New("You must use ONLY one of -c, -t or -x (create, list or extract), plus -f")
	}
	return nil
}

func TarItems(tarfile string, includeFiles []string, options TarOptions) error {
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
		if options.IsAppend {
			//flags = flags | os.O_APPEND
		} else {
			flags = flags | os.O_TRUNC | os.O_CREATE
		}
		zf, err = os.OpenFile(tarfile, flags, mode)
		if options.IsAppend {
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
		out = os.Stdout
	}

	zw := tar.NewWriter(out)
	defer zw.Close()

	for _, itemS := range includeFiles {
		//todo: relative/full path checking
		item := ArchiveItem{itemS, itemS, nil}
		err = addFileToTar(zw, item, options.IsVerbose)
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

func UntarItems(tarfile, destDir string, includeFiles []string, options TarOptions) error {
	var in io.Reader
	if tarfile != "" {
		r, err := os.Open(tarfile)
		if err != nil {
			return err
		}
		defer r.Close()
		in = r
	} else {
		in = os.Stdin
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
			if options.IsVerbose {
				fmt.Printf("Making dir %s:\n", hdr.Name)
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
			if options.IsVerbose {
				fmt.Printf("%s\n", hdr.Name)
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

func TestTarItems(tarfile string, includeFiles []string) error {
	var in io.Reader
	if tarfile != "" {
		r, err := os.Open(tarfile)
		if err != nil {
			return err
		}
		defer r.Close()
		in = r
	} else {
		in = os.Stdin
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
		fmt.Printf("%s %s %s %s\n", modeString, sizeString, modTimeString, hdr.Name)
	}
	return nil
}

func addFileToTar(zw *tar.Writer, item ArchiveItem, isVerbose bool) error {
	if isVerbose {
		fmt.Printf("Adding %s\n", item.FileSystemPath)
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
			err = addFileToTar(zw, ArchiveItem{filepath.Join(item.FileSystemPath, fi.Name()), filepath.Join(item.ArchivePath, fi.Name()), nil}, isVerbose)
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
