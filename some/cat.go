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
	someutils.RegisterPipable(func() someutils.CliPipable { return new(SomeCat) })
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
func (cat *SomeCat) ParseFlags(call []string, errPipe io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("cat", "[options] [files...]", someutils.VERSION)
	flagSet.SetOutput(errPipe)
	flagSet.AliasedBoolVar(&cat.IsShowEnds, []string{"E", "show-ends"}, false, "display $ at end of each line")
	flagSet.AliasedBoolVar(&cat.IsNumber, []string{"n", "number"}, false, "number all output lines")
	flagSet.AliasedBoolVar(&cat.IsSqueezeBlank, []string{"s", "squeeze-blank"}, false, "squeeze repeated empty output lines")

	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}

	if len(flagSet.Args()) > 0 {
		cat.FileNames = flagSet.Args()
	}
	// else it's coming from STDIN
	return nil, 0
}

func (cat *SomeCat) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.AutoPipeErrInOut()
	invocation.AutoHandleSignals()
	if len(cat.FileNames) > 0 {
		for _, fileName := range cat.FileNames {
			file, err := os.Open(fileName)
			if err != nil {
				return err, 1
			} else {
				if cat.isStraightCopy() {
					_, err = io.Copy(invocation.OutPipe, file)
					if err != nil {
						return err, 1
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
							fmt.Fprintf(invocation.OutPipe, "%s%s%s\n", prefix, text, suffix)
						}
						line++
					}
					err := scanner.Err()
					if err != nil {
						return err, 1
					}
				}
				file.Close()
			}
		}
	} else {
		_, err := io.Copy(invocation.OutPipe, invocation.InPipe)
		if err != nil {
			return err, 1
		}
	}
	return nil, 0
}

func NewCat() *SomeCat {
	return new(SomeCat)
}

func Cat(fileNames ...string) someutils.CliPipable {
	cat := NewCat()
	cat.FileNames = fileNames
	return (cat)
}

func CatCli(call []string) (error, int) {
	util := new(SomeCat)
	return someutils.StdInvoke((util), call)
}
