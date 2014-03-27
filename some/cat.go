package some

import (
	"bufio"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"os"
	"strings"
)

func init() {
	someutils.RegisterPipable(func() someutils.PipableCliUtil { return NewCat() })
}

// SomeCat represents and performs a `cat` invocation
type SomeCat struct {
	IsShowEnds     bool
	IsNumber       bool
	IsSqueezeBlank bool
	FileNames      []string
}

func (cat *SomeCat) isStraightCopy() bool {
	return !cat.IsShowEnds && !cat.IsNumber && !cat.IsSqueezeBlank
}

func (cat *SomeCat) Name() string {
	return "cat"
}

func (cat *SomeCat) Number() *SomeCat {
	cat.IsNumber = true
	return cat
}
func (cat *SomeCat) ShowEnds() *SomeCat {
	cat.IsShowEnds = true
	return cat
}
func (cat *SomeCat) SqueezeBlank() *SomeCat {
	cat.IsSqueezeBlank = true
	return cat
}
func (cat *SomeCat) ParseFlags(call []string, errPipe io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("cat", "[options] [files...]", someutils.VERSION)
	flagSet.SetOutput(errPipe)
	flagSet.AliasedBoolVar(&cat.IsShowEnds, []string{"E", "show-ends"}, false, "display $ at end of each line")
	flagSet.AliasedBoolVar(&cat.IsNumber, []string{"n", "number"}, false, "number all output lines")
	flagSet.AliasedBoolVar(&cat.IsSqueezeBlank, []string{"s", "squeeze-blank"}, false, "squeeze repeated empty output lines")

	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errPipe, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	if len(flagSet.Args()) > 0 {
		cat.FileNames = flagSet.Args()
	}
	return nil
}

func (cat *SomeCat) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	if len(cat.FileNames) > 0 {
		for _, fileName := range cat.FileNames {
			if file, err := os.Open(fileName); err == nil {
				if cat.isStraightCopy() {
					_, err = io.Copy(outPipe, file)
					if err != nil {
						return err
					}
				} else {
					scanner := bufio.NewScanner(file)
					line := 1
					var prefix string
					var suffix string
					for scanner.Scan() {
						text := scanner.Text()
						if !cat.IsSqueezeBlank || len(strings.TrimSpace(text)) > 0 {
							if cat.IsNumber {
								prefix = fmt.Sprintf("%d ", line)
							} else {
								prefix = ""
							}
							if cat.IsShowEnds {
								suffix = "$"
							} else {
								suffix = ""
							}
							fmt.Fprintf(outPipe, "%s%s%s\n", prefix, text, suffix)
						}
						line++
					}
					err := scanner.Err()
					if err != nil {
						return err
					}
				}
				file.Close()
			} else {
				return err
			}
		}
	} else {
		_, err := io.Copy(outPipe, inPipe)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewCat() *SomeCat {
	return new(SomeCat)
}

func Cat(fileNames ...string) *SomeCat {
	cat := NewCat()
	cat.FileNames = fileNames
	return cat
}

func CatCli(call []string) error {
	cat := NewCat()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := cat.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return cat.Exec(inPipe, outPipe, errPipe)
}
