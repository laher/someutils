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
	someutils.RegisterPipable(func() someutils.PipableCliUtil { return NewTr() })
}

type SomeTr struct {
	IsDelete     bool
	IsComplement bool
	IsReplace    bool
	set1         string
	set2         string
	inputs       []*regexp.Regexp
	outputs      []string
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
		err = tr.SetSet1(sets[0])
		if err != nil {
			return err
		}
	} else {
		return errors.New("Not enough args supplied")
	}
	if len(sets) > 1 {
		err = tr.SetSet2(sets[1])
		if err != nil {
			return err
		}
	} else if !tr.IsDelete && !tr.IsComplement {
		return errors.New("Not enough args supplied")
	}
	return nil
}

func (tr *SomeTr) SetSet1(set1 string) error {
	tr.set1 = set1
	//tr.inputs, err := convertSet1(tr.set1)
	inputs, err := convertSet1(set1)
	tr.inputs = inputs
	return err
}
func (tr *SomeTr) SetSet2(set2 string) error {
	tr.set2 = set2
	//tr.outputs, err := convertSet2(tr.set2)
	outputs, err := convertSet2(set2)
	tr.outputs = outputs
	return err
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
	/*
		inputs, err := convertSet1(tr.Set1)
		if err != nil {
			return err
		}
		outputs, err := convertSet2(tr.Set2)
		if err != nil {
			return err
		}
	*/
	fu := func(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer, line []byte) error {
		out := line
		for i, reg := range tr.inputs {
			var output string
			if len(tr.outputs) > i {
				output = tr.outputs[i]
			} else {
				output = tr.outputs[len(tr.outputs)-1]
			}
			out = reg.ReplaceAll(out, []byte(output))
		}
		//	fmt.Printf("From %v to %v\n", string(line), string(out))
		_, err := fmt.Fprintln(outPipe, string(out))
		return err
	}
	//fmt.Printf("from %v\n", tr.inputs)
	//fmt.Printf("to %v\n", tr.outputs)

	return someutils.LineProcessor(inPipe, outPipe, errPipe, fu)
}

func NewTr() *SomeTr {
	return new(SomeTr)
}
func Tr(set1, set2 string) *SomeTr {
	tr := NewTr()
	tr.SetSet1(set1)
	tr.SetSet2(set2)
	return tr
}
func TrD(set1 string) *SomeTr {
	tr := NewTr()
	tr.IsDelete = true
	tr.SetSet1(set1)
	return tr
}
func TrC(set1 string) *SomeTr {
	tr := NewTr()
	tr.IsComplement = true
	tr.SetSet1(set1)
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
