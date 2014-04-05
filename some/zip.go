package some

import (
	"archive/zip"
	"errors"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
)

func init() {
	someutils.RegisterPipable(func() someutils.NamedPipable { return NewZip() })
}

// SomeZip represents and performs a `zip` invocation
type SomeZip struct {
	zipFilename string
	items       []string
}

// Name() returns the name of the util
func (z *SomeZip) Name() string {
	return "zip"
}

// ParseFlags parses flags from a commandline []string
func (z *SomeZip) ParseFlags(call []string, errWriter io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("zip", "[options] [files...]", someutils.VERSION)
	flagSet.SetOutput(errWriter)

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}
	args := flagSet.Args()
	if len(args) < 2 {
		flagSet.Usage()
		return errors.New("Not enough args given"), 1
	}
	z.zipFilename = args[0]
	z.items = args[1:]
	return nil, 0
}

// Exec actually performs the zip
func (z *SomeZip) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) (error, int) {
	err := ZipItems(z.zipFilename, z.items)
	if err != nil {
		return err, 1
	}
	return nil, 0

}

func ZipItems(zipFilename string, itemsToArchive []string) error {
	_, err := os.Stat(zipFilename)
	var zf *os.File
	if err != nil {
		if os.IsNotExist(err) {
			zf, err = os.Create(zipFilename)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		zf, err = os.Create(zipFilename)
		if err != nil {
			return err
		}
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	//resources
	for _, itemS := range itemsToArchive {
		//todo: relative/full path checking
		item := someutils.ArchiveItem{itemS, itemS, nil}
		err = addFileToZIP(zw, item)
		if err != nil {
			return err
		}
	}
	//get error where possible
	err = zw.Close()
	return err
}

func addFileToZIP(zw *zip.Writer, item someutils.ArchiveItem) error {
	//fmt.Printf("Adding %s\n", item.FileSystemPath)
	binfo, err := os.Stat(item.FileSystemPath)
	if err != nil {
		return err
	}
	if binfo.IsDir() {
		header, err := zip.FileInfoHeader(binfo)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate
		header.Name = item.ArchivePath
		_, err = zw.CreateHeader(header)
		if err != nil {
			return err
		}
		file, err := os.Open(item.FileSystemPath)
		if err != nil {
			return err
		}
		fis, err := file.Readdir(0)
		for _, fi := range fis {
			err = addFileToZIP(zw, someutils.ArchiveItem{filepath.Join(item.FileSystemPath, fi.Name()), filepath.Join(item.ArchivePath, fi.Name()), nil})
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(binfo)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate
		header.Name = item.ArchivePath
		w, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		bf, err := os.Open(item.FileSystemPath)
		if err != nil {
			return err
		}
		defer bf.Close()
		_, err = io.Copy(w, bf)
		if err != nil {
			return err
		}
	}
	return err
}

// Factory for *SomeZip
func NewZip() *SomeZip {
	return new(SomeZip)
}

// Factory for *SomeZip
func Zip(args ...string) *SomeZip {
	z := NewZip()
	z.zipFilename = args[0]
	z.items = args[1:]
	return z
}

// CLI invocation for *SomeZip
func ZipCli(call []string) (error, int) {
	z := NewZip()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err, code := z.ParseFlags(call, errPipe)
	if err != nil {
		return err, code
	}
	return z.Exec(inPipe, outPipe, errPipe)
}
