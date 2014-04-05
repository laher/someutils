package some

import (
	"bytes"
	"github.com/laher/someutils"
	"io"
	"strings"
	"testing"
	"time"
)

func TestTrCli(t *testing.T) {
	//var out bytes.Buffer
	//var errout bytes.Buffer
	//inPipe, outPipe, errPipe := (strings.NewReader("HI"), &out, &errout)
	err, code := TrCli([]string{"tr", "I", "O"})
	if err != nil {
		if code != 0 {
			t.Errorf("Error: %v, code: %d\n", err, code)
		}
		t.Logf("Error: %v, code: %d\n", err, code)
	}
	//println(out.String())
}

func TestFluentTr(t *testing.T) {
	var out bytes.Buffer
	var errout bytes.Buffer
	inPipe, outPipe, errPipe := strings.NewReader("HI"), &out, &errout
	err, code := Tr("I", "O").Exec(inPipe, outPipe, errPipe)
	if err != nil {
		if code != 0 {
			t.Errorf("Error: %v, code: %d\n", err, code)
		}
		t.Logf("Error: %v, code: %d\n", err, code)
	}
	t.Log(out.String())
}

func Test2pipes(t *testing.T) {
	var out bytes.Buffer
	var errout bytes.Buffer
	r, w := io.Pipe()
	in := strings.NewReader("Hi\nHo\nhI\nhO\n")
	inPipe1, outPipe1, errPipe1 := in, w, &errout
	inPipe2, outPipe2, errPipe2 := r, &out, &errout
	tr1 := Tr("H", "O")
	tr2 := Tr("I", "J")
	go tr1.Exec(inPipe1, outPipe1, errPipe1)
	go tr2.Exec(inPipe2, outPipe2, errPipe2)
	time.Sleep(1 * time.Second)
	output := out.String()
	expected := "Oi\nOo\nhJ\nhO\n"
	if output != expected {
		t.Error("Expected ", expected, ", got ", output)
	}
	t.Logf("Errout: %+v\n", errout.String())
}

func TestTrPipeline(t *testing.T) {
	var out bytes.Buffer
	var errout bytes.Buffer
	in := strings.NewReader("Hi\nHo\nhI\nhO\n")
	pipeline := someutils.NewPipeline(Tr("H", "O"), Tr("I", "J"))
	e := pipeline.Exec(someutils.NewPipeset(in, &out, &errout))
	err, code, index := someutils.Wait(e, 2)
	if err != nil {
		t.Logf("Errout: %+v\n", errout.String())
		if code != 0 {
			t.Errorf("Error: %v, code: %d, index: %d\n", err, code, index)
		}
		t.Logf("Error: %v, code: %d\n", err, code)
	}

	output := out.String()
	expected := "Oi\nOo\nhJ\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}
