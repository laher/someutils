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
	someutils.RegisterSimple(func() someutils.CliPipableSimple { return new(SomeTr) })
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
func (tr *SomeTr) ParseFlags(call []string, errWriter io.Writer) (error, int) {
	flagSet := uggo.NewFlagSetDefault("tr", "[OPTION]... SET1 [SET2]", someutils.VERSION)
	flagSet.SetOutput(errWriter)
	flagSet.AliasedBoolVar(&tr.IsDelete, []string{"d", "delete"}, false, "Delete characters in SET1, do not translate")
	err, code := flagSet.ParsePlus(call[1:])
	if err != nil {
		return err, code
	}
	sets := flagSet.Args()
	if len(sets) > 0 {
		err = tr.SetSet1(sets[0])
		if err != nil {
			return err, 1
		}
	} else {
		return errors.New("Not enough args supplied"), 1
	}
	if len(sets) > 1 {
		err = tr.SetSet2(sets[1])
		if err != nil {
			return err, 1
		}
	} else if !tr.IsDelete && !tr.IsComplement {
		return errors.New("Not enough args supplied"), 1
	}
	return nil, 0
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

//TODO fix behaviour of set1/set2 relationship
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

func (tr *SomeTr) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) (error, int) {
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
	fu := func(inPipe io.Reader, outPipe2 io.Writer, errPipe io.Writer, line []byte) error {
		//println("tr processing line")
		out := string(line)
		for i, reg := range tr.inputs {
			var output string
			if len(tr.outputs) > i {
				output = tr.outputs[i]
			} else {
				output = tr.outputs[len(tr.outputs)-1]
			}
			out = reg.ReplaceAllString(out, output)
		}
		//	fmt.Printf("From %v to %v\n", string(line), string(out))
		//println("tr printing line: |", out, "|")
		_, err := fmt.Fprintln(outPipe2, out)
		//println("tr processed line")
		return err
	}
	//fmt.Printf("from %v\n", tr.inputs)
	//fmt.Printf("to %v\n", tr.outputs)

	err := someutils.LineProcessor(inPipe, outPipe, errPipe, fu)
	if err != nil {
		return err, 1
	}
	return nil, 0
}

func NewTr() *SomeTr {
	return new(SomeTr)
}
func Tr(set1, set2 string) someutils.NamedPipable {
	tr := NewTr()
	tr.SetSet1(set1)
	tr.SetSet2(set2)
	return someutils.WrapNamed(tr)
}
func TrD(set1 string) someutils.NamedPipable {
	tr := NewTr()
	tr.IsDelete = true
	tr.SetSet1(set1)
	return someutils.WrapNamed(tr)
}
func TrC(set1 string) someutils.NamedPipable {
	tr := NewTr()
	tr.IsComplement = true
	tr.SetSet1(set1)
	return someutils.WrapNamed(tr)
}

func TrCli(call []string) (error, int) {

	util := new(SomeTr)
	return someutils.StdInvoke(someutils.WrapUtil(util), call)
}
