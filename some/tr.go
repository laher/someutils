package some

import (
	"errors"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"regexp"
	"strings"
)

func init() {
	someutils.RegisterSome(func() someutils.SomeUtil { return NewTr() })
}

type SomeTr struct {
	IsDelete     bool
	IsComplement bool
	IsReplace    bool
	Set1         string
	Set2         string
}

func (tr *SomeTr) Name() string {
	return "tr"
}
func (tr *SomeTr) ParseFlags(call []string, errWriter io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("tr", "[OPTION]... SET1 [SET2]", someutils.VERSION)
	flagSet.SetOutput(errWriter)
	flagSet.AliasedBoolVar(&tr.IsDelete, []string{"d", "delete"}, false, "Delete characters in SET1, do not translate")
	err := flagSet.Parse(call[1:])
	if err != nil {
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	sets := flagSet.Args()
	if len(sets) > 0 {
		tr.Set1 = sets[0]
	} else {
		return errors.New("Not enough args supplied")
	}
	if len(sets) > 1 {
		tr.Set2 = sets[1]
	} else if !tr.IsDelete && !tr.IsComplement {
		return errors.New("Not enough args supplied")
	}
	return nil
}

func convertSet2(set2 string) ([]string, error) {
	if strings.Contains(set2, "-") {
		parts := strings.Split(set2, "-")
		firstChar := parts[0]
		unicodePointStart := int(firstChar[0])
		lastChar := parts[1]
		unicodePointEnd := int(lastChar[0])
		outputs := []string{}
		for i := unicodePointStart; i <= unicodePointEnd; i++ {
			st := string([]byte{byte(i)})
			outputs = append(outputs, st)
		}
		return outputs, nil
	} else {
		return []string{set2}, nil
	}
}
func convertSet1(set1 string) ([]*regexp.Regexp, error) {
	inputs := []*regexp.Regexp{}
	if strings.Contains(set1, "-") {
		parts := strings.Split(set1, "-")
		firstChar := parts[0]
		unicodePointStart := int(firstChar[0])
		lastChar := parts[1]
		unicodePointEnd := int(lastChar[0])
		for i := unicodePointStart; i <= unicodePointEnd; i++ {
			r, err := regexp.Compile(string([]byte{byte(i)}))
			if err != nil {
				return inputs, err
			}
			inputs = append(inputs, r)
		}
	} else {
		r, err := regexp.Compile("[" + set1 + "]")
		if err != nil {
			return inputs, err
		}
		inputs = append(inputs, r)
	}
	return inputs, nil
}

func (tr *SomeTr) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	inputs, err := convertSet1(tr.Set1)
	if err != nil {
		return err
	}
	outputs, err := convertSet2(tr.Set2)
	if err != nil {
		return err
	}
	fu := func(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer, line []byte) error {
		out := line
		for i, reg := range inputs {
			var output string
			if len(outputs) > i {
				output = outputs[i]
			} else {
				output = outputs[len(outputs)-1]
			}
			out = reg.ReplaceAll(out, []byte(output))
		}
		_, err = fmt.Fprintln(outPipe, string(out))
		return err
	}
	return someutils.LineProcessor(inPipe, outPipe, errPipe, fu)
}

func NewTr() *SomeTr {
	return new(SomeTr)
}
func Tr(set1, set2 string) *SomeTr {
	tr := NewTr()
	tr.Set1 = set1
	tr.Set2 = set2
	return tr
}
func TrD(set1 string) *SomeTr {
	tr := NewTr()
	tr.IsDelete = true
	tr.Set1 = set1
	return tr
}
func TrC(set1 string) *SomeTr {
	tr := NewTr()
	tr.IsComplement = true
	tr.Set1 = set1
	return tr
}

func TrCli(call []string) error {
	tr := NewTr()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := tr.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return tr.Exec(inPipe, outPipe, errPipe)
}
