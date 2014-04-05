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
	someutils.RegisterPipable(func() someutils.NamedPipable { return new(SomeGunzip) })
}

// SomeGunzip represents and performs a `gunzip` invocation
type SomeGunzip struct {
	IsTest    bool
	IsKeep    bool
	IsPipeOut bool
	Filenames []string
}

// Name() returns the name of the util
func (gunzip *SomeGunzip) Name() string {
	return "gunzip"
}

// ParseFlags parses flags from a commandline []string
func (gunzip *SomeGunzip) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("gunzip", "[options] file.gz [list...]", someutils.VERSION)
	flagSet.SetOutput(errPipe)
	flagSet.AliasedBoolVar(&gunzip.IsTest, []string{"t", "test"}, false, "test archive data")
	flagSet.AliasedBoolVar(&gunzip.IsKeep, []string{"k", "keep"}, false, "keep gzip file")
	flagSet.AliasedBoolVar(&gunzip.IsPipeOut, []string{"c", "stdout", "is-stdout"}, false, "output will go to the standard output")

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}
	args := flagSet.Args()
	//TODO STDIN support
	if len(args) > 0 {
		//OK
	} else if uggo.IsPipingStdin() {
		//OK
	} else {
		return errors.New("No gzip filename given"), 1
	}
	gunzip.Filenames = args
	return nil, 0
}

// Exec actually performs the gunzip
func (gunzip *SomeGunzip) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) (error, int) {
	if gunzip.IsTest {
		err := TestGzipItems(gunzip.Filenames)
		if err != nil {
			return err, 1
		}
	} else {
		err := gunzip.gunzipItems(inPipe, outPipe, errPipe)
		if err != nil {
			return err, 1
		}
	}
	return nil, 0

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

func (gunzip *SomeGunzip) gunzipItems(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	if len(gunzip.Filenames) == 0 {
		//stdin
		err := gunzip.gunzipItem(inPipe, outPipe, errPipe)
		if err != nil {
			return err
		}
	} else {
		for _, item := range gunzip.Filenames {
			fh, err := os.Open(item)
			if err != nil {
				return err
			}
			err = gunzip.gunzipItem(fh, outPipe, errPipe)
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
	}
	return nil
}

func (gunzip *SomeGunzip) gunzipItem(item io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	r, err := gzip.NewReader(item)
	if err != nil {
		return err
	}
	defer r.Close()
	var w io.Writer
	if gunzip.IsPipeOut {
		w = outPipe
		_, err = io.Copy(w, r)
		if err != nil {
			return err
		}
	} else {
		destFileName := r.Header.Name
		fmt.Fprintln(errPipe, "Filename", destFileName)
		destFile, err := os.Create(destFileName)
		defer destFile.Close()
		if err != nil {
			return err
		}
		w = destFile
		_, err = io.Copy(w, r)
		if err != nil {
			return err
		}

		err = destFile.Close()
		if err != nil {
			return err
		}
	}
	err = r.Close()
	return err
}


func GunzipToOut(args ...string) *SomeGunzip {
	gunzip := Gunzip(args...)
	gunzip.IsPipeOut = true
	return gunzip
}

// Factory for *SomeGunzip
func Gunzip(args ...string) *SomeGunzip {
	gunzip := new(SomeGunzip)
	gunzip.Filenames = args
	if len(args) == 0 {
		gunzip.IsPipeOut = true
	}
	return gunzip
}

// CLI invocation for *SomeGunzip
func GunzipCli(call []string) (error, int) {
	gunzip := new(SomeGunzip)
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err, code := gunzip.ParseFlags(call, errPipe)
	if err != nil {
		return err, code
	}
	return gunzip.Exec(inPipe, outPipe, errPipe)
}
