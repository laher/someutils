package someutils

import (
	"compress/gzip"
	"errors"
	"github.com/laher/uggo"
	"io"
	"os"
)

type GunzipOptions struct {
	IsTest bool
	IsKeep bool
}

func init() {
	Register(Util{
		"gunzip",
		Gunzip})
}

func Gunzip(call []string) error {
	options := GunzipOptions{}
	flagSet := uggo.NewFlagSetDefault("gunzip", "[options] file.gz [list]", VERSION)
	flagSet.AliasedBoolVar(&options.IsTest, []string{"t", "test"}, false, "test archive data")
	flagSet.AliasedBoolVar(&options.IsKeep, []string{"k", "keep"}, false, "keep gzip file")

	err := flagSet.Parse(call[1:])
	if err != nil {
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
	if options.IsTest {
		err = TestGzipItems(args)
		if err != nil {
			return err
		}
	} else {
		err = GunzipItems(args, options)
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

func GunzipItems(items []string, options GunzipOptions) error {
	for _, item := range items {
		fh, err := os.Open(item)
		if err != nil {
			return err
		}
		err = GunzipItem(fh)
		if err != nil {
			return err
		}
		err = fh.Close()
		if err != nil {
			return err
		}
		if !options.IsKeep {
			err = os.Remove(item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GunzipItem(item io.Reader) error {
	r, err := gzip.NewReader(item)
	if err != nil {
		return err
	}
	defer r.Close()
	destFileName := r.Header.Name
	println("Filename", destFileName)
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
