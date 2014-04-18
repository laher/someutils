package some

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	//fmt.Printf("invocation: %+v\n", invocation)
	//fmt.Printf("gz: %+v\n", gz)
	invocation.ErrPipe.Drain()
	invocation.AutoHandleSignals()
	if len(gz.Filenames) < 1 {
		//pipe in?
		var writer io.Writer
		var outputFilename string
		if gz.outFile != "" {
			outputFilename = gz.outFile
			var err error
			writer, err = os.Create(outputFilename)
			if err != nil {
				return err, 1
			}
		} else {
		//	fmt.Printf("stdin to stdout: %+v\n", gz)
			outputFilename = "S" //seems to be the default used by gzip
			writer = invocation.MainPipe.Out
		}
		err := gz.doGzip(invocation.MainPipe.In, writer, filepath.Base(outputFilename))
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
				writer = invocation.MainPipe.Out
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

	inw := new (bytes.Buffer)
	_, err := io.Copy(inw, reader)
	if err != nil {
		return err
	}
//	fmt.Println("input: ", inw.String())
	rdr := strings.NewReader(inw.String())

	outw := new (bytes.Buffer)
	


	gzw := gzip.NewWriter(outw)
	defer gzw.Close()
	gzw.Header.Comment = "Gzipped by someutils."
	gzw.Header.Name = filename

//	fmt.Println("Copying ")
	_, err = io.Copy(gzw, rdr)
//	fmt.Println("Copied ", i)
	if err != nil {
		fmt.Println("Copied err", err)
		return err
	}
	//get error where possible
	err = gzw.Close()
	//fmt.Println("Closed ")
	if err != nil {
		fmt.Println("Closed err", err)
		return err
	}
//	fmt.Println("Wrote OK ", i)
	
//	fmt.Println("Wrote: ", outw.String())
	_, err = io.Copy(writer, outw)
//	fmt.Printf("Copied to eventual writer %T. length: %d\n", outw, i)
	
	return err
}

// Factory for *SomeGzip
func Gzip(args ...string) someutils.CliPipable {
	gz := new(SomeGzip)
	gz.Filenames = args
	if len(args) < 1 {
		gz.IsStdout = true
	}
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
