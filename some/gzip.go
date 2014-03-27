package some

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
)

func init() {
	someutils.RegisterPipable(func() someutils.PipableCliUtil { return NewGzip() })
}

// SomeGzip represents and performs a `gzip` invocation
type SomeGzip struct {
	IsKeep bool

	Filenames []string
}

// Name() returns the name of the util
func (gz *SomeGzip) Name() string {
	return "gzip"
}

// TODO: add validation here

// ParseFlags parses flags from a commandline []string
func (gz *SomeGzip) ParseFlags(call []string, errPipe io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("gzip", "[options] [files...]", someutils.VERSION)
	flagSet.SetOutput(errPipe)

	flagSet.AliasedBoolVar(&gz.IsKeep, []string{"k", "keep"}, false, "keep gzip file")

	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errPipe, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}

	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	args := flagSet.Args()
	//TODO STDIN support
	if len(args) < 1 {
		flagSet.Usage()
		return errors.New("Not enough args given")
	}
	gz.Filenames = args
	return nil
}

// Exec actually performs the gzip
func (gz *SomeGzip) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	return GzipItems(gz.Filenames, gz)
}

func GzipItems(itemsToCompress []string, gz *SomeGzip) error {
	for _, item := range itemsToCompress {
		err := GzipItem(item)
		if err != nil {
			return err
		}
		if !gz.IsKeep {
			err = os.Remove(item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GzipItem(filename string) error {
	gzipItem, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer gzipItem.Close()
	//todo use tgz for tars?
	gzipFilename := filename + ".gz"
	gzf, err := os.Create(gzipFilename)
	if err != nil {
		return err
	}
	defer gzf.Close()

	gzw := gzip.NewWriter(gzf)
	defer gzw.Close()
	gzw.Header.Comment = "file compressed by someutils-gzip"
	gzw.Header.Name = filepath.Base(filename)

	_, err = io.Copy(gzw, gzipItem)
	if err != nil {
		return err
	}
	//get error where possible
	err = gzw.Close()
	if err != nil {
		return err
	}

	return nil
}

// Factory for *SomeGzip
func NewGzip() *SomeGzip {
	return new(SomeGzip)
}

// Fluent factory for *SomeGzip
func Gzip(args ...string) *SomeGzip {
	gz := NewGzip()
	gz.Filenames = args
	return gz
}

// CLI invocation for *SomeGzip
func GzipCli(call []string) error {
	gz := NewGzip()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := gz.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return gz.Exec(inPipe, outPipe, errPipe)
}
