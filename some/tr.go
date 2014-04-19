package some

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"regexp"
	"strings"
)

func init() {
	someutils.RegisterPipable(func() someutils.CliPipable { return new(SomeTr) })
}

// SomeTr mimics the functionality of the `tr` function from coreutils.
type SomeTr struct {
	IsDelete     bool
	IsSqueeze    bool
	IsComplement bool // currently unused as I just don't get it.
	set1         string
	set2         string
	translations map[*regexp.Regexp]string
}

// returns the name of the command
func (tr *SomeTr) Name() string {
	return "tr"
}

// ParseFlags parses commandline options.
func (tr *SomeTr) ParseFlags(call []string, errWriter io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("tr", "[OPTION]... SET1 [SET2]", someutils.VERSION)
	flagSet.SetOutput(errWriter)
	flagSet.AliasedBoolVar(&tr.IsDelete, []string{"d", "delete"}, false, "Delete characters in SET1, do not translate")
	flagSet.AliasedBoolVar(&tr.IsSqueeze, []string{"s", "squeeze-repeats"}, false, "replace each input sequence of a repeated character that is listed in SET1 with a single occurence of that character")
	//Don't get the Complement thing. Just don't understand it right now.
	flagSet.AliasedBoolVar(&tr.IsComplement, []string{"c", "complement"}, false, "use the complement of SET1")
	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}
	sets := flagSet.Args()
	if len(sets) > 0 {
		tr.set1 = sets[0]
	} else {
		return errors.New("Not enough args supplied"), 1
	}
	if len(sets) > 1 {
		tr.set2 = sets[1]
	} else if !tr.IsDelete {
		return errors.New("Not enough args supplied"), 1
	}
	err = tr.Preprocess()
	if err != nil {
		return err, 1
	}
	return nil, 0
}

// Parses SET1 and SET2 into a map of translations
func (tr *SomeTr) Preprocess() error {
	tr.translations = map[*regexp.Regexp]string{}
	set1 := tr.set1

	var set1Part, set2Part string
	var set2len int
	var err error
	if tr.IsDelete {
		// processSet1 only
		for len(set1) > 0 {
			set1Part, set1, _ = nextPartSet1(set1)
			reg, err := tr.toRegexp(set1Part)
			if err != nil {
				return err
			}
			//replacement is empty-string
			tr.translations[reg] = ""
		}
	} else {
		// process both sets together
		set2 := tr.set2
		for len(set1) > 0 {
			set1Part, set1, set2len = nextPartSet1(set1)
			var set2New string
			set2Part, set2New, err = nextPartSet2(set2, set2len)
			if len(set2New) > 0 { //incase set2 is shorter than set1, behave like BSD tr (rather than SystemV, which truncates set1 instead)
				set2 = set2New
			}
			if err != nil {
				return err
			}
			reg, err := tr.toRegexp(set1Part)
			if err != nil {
				return err
			}
			tr.translations[reg] = set2Part
		}
	}
	return nil
}

func (tr *SomeTr) toRegexp(set1Part string) (*regexp.Regexp, error) {
	maybeSqueeze := ""
	maybeComplement := ""
	if tr.IsSqueeze {
		maybeSqueeze = "+"
	}
	if tr.IsComplement {
		maybeComplement = "^"
	}
	regString := "^["+maybeComplement+set1Part+"]"+maybeSqueeze
	//fmt.Println(regString)
	reg, err := regexp.Compile(regString)
	return reg, err
}

// Parser for SET1. Supports single-chars, ranges, and regex-like character groups
func nextPartSet1(set1 string) (string, string, int) {
	if strings.HasPrefix(set1, "[") {
		//find matching
		if strings.Contains(set1, "]") {
			return set1[:strings.Index(set1, "]")+1], set1[strings.Index(set1, "]")+1:], 1

		} else {
			return set1[:1], set1[1:], 1
		}
	} else if len(set1)>2 && set1[1]=='-' {
		return set1[:3], set1[3:], 1
	} else {
		return set1[:1], set1[1:], 1
	}
}

// Parser for SET2. Supports single and multiple chars
func nextPartSet2(set2 string, set2len int) (string, string, error) {
	if len(set2) < set2len {
		return "", "", errors.New(fmt.Sprintf("Error out of range (%d - %s)", set2len, set2))
	}
	return set2[:set2len], set2[set2len:], nil
}

// Invoke actually carries out the command
func (tr *SomeTr) Invoke(invocation *someutils.Invocation) (error, int) {
	invocation.ErrPipe.Drain()
	invocation.AutoHandleSignals()
	tr.Preprocess()
	fu := func(inPipe io.Reader, outPipe2 io.Writer, errPipe io.Writer, line []byte) error {
		//println("tr processing line")
		//outline := line
		var buffer bytes.Buffer
		remainder := string(line)
		for len(remainder) > 0 {
			trimLeft := 1
			nextPart := remainder[:trimLeft]
			for reg, v := range tr.translations {
				//fmt.Printf("Translation '%v'=>'%s' on '%s'\n", reg, v, remainder)
				match := reg.MatchString(remainder)
				if match {
					toReplace := reg.FindString(remainder)
					replacement := reg.ReplaceAllString(toReplace, v)
					//fmt.Printf("Replace %s=>%s\n", toReplace, replacement)
					nextPart = replacement
					//if squeezing has taken place, remove more leading chars accordingly
					trimLeft = len(toReplace)
					if !tr.IsComplement {
						break
					}
				} else if tr.IsComplement {
					// this is a double-negative - non-match of negative-regex.
					// This implies that set1 matches the current input character.
					// So, keep it as-is and break out of the loop.
					trimLeft = 1
					nextPart = remainder[:trimLeft]
					break
				}
			}

			remainder = remainder[trimLeft:]
			buffer.WriteString(nextPart)
		}
		out := buffer.String()
		_, err := fmt.Fprintln(outPipe2, out)
		return err
	}
	err := someutils.LineProcessor(invocation.MainPipe.In, invocation.MainPipe.Out, invocation.ErrPipe.Out, fu)
	if err != nil {
		return err, 1
	}
	return nil, 0
}

func NewTr() *SomeTr {
	return new(SomeTr)
}

// Factory for normal `tr`
func Tr(set1, set2 string) someutils.CliPipable {
	tr := NewTr()
	tr.set1 = set1
	tr.set2 = set2
	return (tr)
}

// Factory for `tr -d` (delete set1)
func TrD(set1 string) someutils.CliPipable {
	tr := NewTr()
	tr.IsDelete = true
	tr.set1 = set1
	return (tr)
}

// Factory for `tr -c` (complement)
func TrC(set1 string, set2 string) someutils.CliPipable {
	tr := NewTr()
	tr.IsComplement = true
	tr.set1 = set1
	tr.set2 = set2
	return (tr)
}

//Cli helper for tr
func TrCli(call []string) (error, int) {

	util := new(SomeTr)
	return someutils.StdInvoke((util), call)
}
