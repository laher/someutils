package someutils

import (
	"compress/gzip"
	"errors"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
)

type GzipOptions struct {
	IsKeep bool
}

func init() {
	Register(Util{
		"gzip",
		Gzip})
}

func Gzip(call []string) error {
	options := GzipOptions{}
	flagSet := uggo.NewFlagSetDefault("gzip", "[options] [files...]", VERSION)
	flagSet.AliasedBoolVar(&options.IsKeep, []string{"k", "keep"}, false, "keep gzip file")
	err := flagSet.Parse(call[1:])
	if err != nil {
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
	err = GzipItems(args, options)
	if err != nil {
		return err
	}
	return nil
}

func GzipItems(itemsToCompress []string, options GzipOptions) error {
	for _, item := range itemsToCompress {
		err := GzipItem(item)
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

