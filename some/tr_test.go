package some

import (
	"github.com/laher/someutils"
	"testing"
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
/*
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
*/
func TestTrPipeline(t *testing.T) {
	pipeline := someutils.NewPipeline(Tr("H", "O"), Tr("I", "J"))
	invocation, out, errout := someutils.InvocationFromString("Hi\nHo\nhI\nhO\n")
	err, invocationchan, count := invocation.PipeToPipeline(pipeline)
	if err != nil {
		t.Errorf("error piping to pipeline: %v", err)
	}
	errinvocation := someutils.Wait(invocationchan, count)
	outstring := out.String()
	if errinvocation.Err!=nil {
		t.Logf("errout: %+v\n", errout.String())
		t.Logf("stdout: %+v", outstring)
		t.Logf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			t.Errorf("error: %+v\n", errinvocation.Err)
		}
	}

	output := out.String()
	expected := "Oi\nOo\nhJ\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}
