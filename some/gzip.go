package some

import (
	"compress/gzip"
	"errors"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
)

func init() {
	someutils.RegisterPipable(func() someutils.CliPipable { return new(SomeGzip) })
}

// SomeGzip represents and performs a `gzip` invocation
type SomeGzip struct {
	IsKeep    bool
	IsStdout  bool
	Filenames []string
	outFile   string
}

// Name() returns the name of the util
func (gz *SomeGzip) Name() string {
	return "gzip"
}

// ParseFlags parses flags from a commandline []string
func (gz *SomeGzip) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("gzip", "[options] [files...]", someutils.VERSION)
	flagSet.SetOutput(errPipe)

	flagSet.AliasedBoolVar(&gz.IsKeep, []string{"k", "keep"}, false, "keep gzip file")
	flagSet.AliasedBoolVar(&gz.IsStdout, []string{"c", "stdout", "to-stdout"}, false, "pipe output to standard out. Keep source file.")

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}
	args := flagSet.Args()
	//TODO STDIN support
	if len(args) < 1 {
		flagSet.Usage()
		return errors.New("Not enough args given"), 1
	}
	gz.Filenames = args
	return nil, 0
}

// Exec actually performs the gzip
func (gz *SomeGzip) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.AutoPipeErrInOut()
	invocation.AutoHandleSignals()
	if len(gz.Filenames) == 0 {
		//pipe in?
		var writer io.Writer
		outputFilename := ""
		if gz.outFile != "" {
			outputFilename = gz.outFile
			var err error
			writer, err = os.Create(outputFilename)
			if err != nil {
				return err, 1
			}
		} else {
			outputFilename = ""
			writer = invocation.OutPipe
		}
		err := gz.doGzip(invocation.InPipe, writer, filepath.Base(outputFilename))
		if err != nil {
			return err, 1
		}
	} else {
		//todo make sure it closes saved file cleanly
		for _, inputFilename := range gz.Filenames {
			inputFile, err := os.Open(inputFilename)
			if err != nil {
				return err, 1
			}
			defer inputFile.Close()

			var writer io.Writer
			if !gz.IsStdout {
				outputFilename := inputFilename + ".gz"
				gzf, err := os.Create(outputFilename)
				if err != nil {
					return err, 1
				}
				defer gzf.Close()
				writer = gzf
			} else {
				writer = invocation.OutPipe
			}
			err = gz.doGzip(inputFile, writer, filepath.Base(inputFilename))
			if err != nil {
				return err, 1
			}

			err = inputFile.Close()
			if err != nil {
				return err, 1
			}

			// only remove source if specified and possible
			if !gz.IsKeep && !gz.IsStdout {
				err = os.Remove(inputFilename)
				if err != nil {
					return err, 1
				}
			}
		}
	}
	return nil, 0
}

func (gz *SomeGzip) doGzip(reader io.Reader, writer io.Writer, filename string) error {
	gzw := gzip.NewWriter(writer)
	defer gzw.Close()
	gzw.Header.Comment = "file compressed by someutils-gzip"
	gzw.Header.Name = filename

	_, err := io.Copy(gzw, reader)
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
func Gzip(args ...string) someutils.CliPipable {
	gz := new(SomeGzip)
	gz.Filenames = args
	return (gz)
}

// Factory for *SomeGzip
func GzipTo(outFile string) someutils.CliPipable {
	gz := new(SomeGzip)
	gz.outFile = outFile
	return (gz)
}

// CLI invocation for *SomeGzip
func GzipCli(call []string) (error, int) {
	util := new(SomeGzip)
	return someutils.StdInvoke((util), call)
}
