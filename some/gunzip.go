package some

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
)

func init() {
	someutils.RegisterSome(func() someutils.SomeUtil { return NewGunzip() })
}

// SomeGunzip represents and performs a `gunzip` invocation
type SomeGunzip struct {
	IsTest    bool
	IsKeep    bool
	Filenames []string
}

// Name() returns the name of the util
func (gunzip *SomeGunzip) Name() string {
	return "gunzip"
}

// TODO: add validation here

// ParseFlags parses flags from a commandline []string
func (gunzip *SomeGunzip) ParseFlags(call []string, errWriter io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("gunzip", "[options] file.gz [list...]", someutils.VERSION)
	flagSet.SetOutput(errWriter)
	flagSet.AliasedBoolVar(&gunzip.IsTest, []string{"t", "test"}, false, "test archive data")
	flagSet.AliasedBoolVar(&gunzip.IsKeep, []string{"k", "keep"}, false, "keep gzip file")

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
	//TODO STDIN support
	if len(args) < 1 {
		return errors.New("No gzip filename given")
	}
	gunzip.Filenames = args
	return nil
}

// Exec actually performs the gunzip
func (gunzip *SomeGunzip) Exec(pipes someutils.Pipes) error {
	if gunzip.IsTest {
		err := TestGzipItems(gunzip.Filenames)
		if err != nil {
			return err
		}
	} else {
		err := GunzipItems(gunzip.Filenames, gunzip, pipes.Out())
		if err != nil {
			return err
		}
	}
	return nil

}

func TestGzipItems(items []string) error {
	for _, item := range items {
		fh, err := os.Open(item)
		if err != nil {
			return err
		}
		err = TestGzipItem(fh)
		if err != nil {
			return err
		}
	}
	return nil
}

//TODO: proper file checking
func TestGzipItem(item io.Reader) error {
	r, err := gzip.NewReader(item)
	if err != nil {
		return err
	}
	defer r.Close()
	return nil
}

func GunzipItems(items []string, gunzip *SomeGunzip, outPipe io.Writer) error {
	for _, item := range items {
		fh, err := os.Open(item)
		if err != nil {
			return err
		}
		err = GunzipItem(fh, outPipe)
		if err != nil {
			return err
		}
		err = fh.Close()
		if err != nil {
			return err
		}
		if !gunzip.IsKeep {
			err = os.Remove(item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GunzipItem(item io.Reader, outPipe io.Writer) error {
	r, err := gzip.NewReader(item)
	if err != nil {
		return err
	}
	defer r.Close()
	destFileName := r.Header.Name
	fmt.Fprintln(outPipe, "Filename", destFileName)
	destFile, err := os.Create(destFileName)
	defer destFile.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(destFile, r)
	if err != nil {
		return err
	}
	err = destFile.Close()
	if err != nil {
		return err
	}
	err = r.Close()
	return err
}

// Factory for *SomeGunzip
func NewGunzip() *SomeGunzip {
	return new(SomeGunzip)
}

// Fluent factory for *SomeGunzip
func Gunzip(args ...string) *SomeGunzip {
	gunzip := NewGunzip()
	gunzip.Filenames = args
	return gunzip
}

// CLI invocation for *SomeGunzip
func GunzipCli(call []string) error {
	gunzip := NewGunzip()
	pipes := someutils.StdPipes()
	err := gunzip.ParseFlags(call, pipes.Err())
	if err != nil {
		return err
	}
	return gunzip.Exec(pipes)
}
