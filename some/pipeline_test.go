package some

import (
	"bytes"
	. "github.com/laher/someutils"
//	"io"
	"strings"
	"testing"
	"time"
)

func TestPipeline1(t *testing.T) {
	var out, errout bytes.Buffer
	pipeline := NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O")) //, RedirectTo("file.txt"))
	err, code, index := pipeline.ExecAndWait(NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), &out, &errout))
	if err!=nil {
		t.Logf("Errout: %+v\n", errout.String())
		t.Logf("Stdout: %+v", out.String())
		t.Logf("Error: %+v, code: %d, index: %d\n", err, code, index)
		if code != 0 {
			t.Errorf("Error: %+v\n", err)
		}
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}

func TestPipeline2(t *testing.T) {
	var out, errout bytes.Buffer
	pipeline := NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"))
	err, code, index := pipeline.ExecAndWait(NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), &out, &errout))
	if err!=nil {
		t.Logf("Errout: %+v\n", errout.String())
		t.Logf("Stdout: %+v", out.String())
		t.Logf("Error: %+v, code: %d, index: %d\n", err, code, index)
		if code != 0 {
			t.Errorf("Error: %+v\n", err)
		}
	}

	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}

}

func TestRedirect1(t *testing.T) {
	var out, errout bytes.Buffer
	pipeline := NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"), OutTo("test.txt"), Cat("test.txt"))
	err, code, index := pipeline.ExecAndWait(NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), &out, &errout))
	if err!=nil {
		t.Logf("Errout: %+v\n", errout.String())
		t.Logf("Stdout: %+v", out.String())
		t.Logf("Error: %+v, code: %d, index: %d\n", err, code, index)
		if code != 0 {
			t.Errorf("Error: %+v\n", err)
		}
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}

func TestRedirectOutErr(t *testing.T) {
	var out, errout bytes.Buffer
	pipeline := NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"), OutToErr())
	err, code, index := pipeline.ExecAndWait(NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), &out, &errout))
	if err!=nil {
		t.Logf("Errout: %+v\n", errout.String())
		t.Logf("Stdout: %+v", out.String())
		t.Logf("Error: %+v, code: %d, index: %d\n", err, code, index)
		if code != 0 {
			t.Errorf("Error: %+v\n", err)
		}
	}
	output := out.String()
	outputErr := errout.String()
	expected := "Oi\nOo\nhO\n"
	if outputErr != expected {
		t.Logf("Errout: %+v", outputErr)
		t.Logf("Stdout: %+v", output)
		t.Error("Expected\n ", expected, ", Got:\n ", outputErr)
	}
}
func TestRedirectOutNull(t *testing.T) {
	out := new(bytes.Buffer)
	errout := new (bytes.Buffer)
	in := bytes.NewBufferString("Hi\nHo\nhI\nhO\n")
/*
	inR, inW := io.Pipe()
	closef := func () {
		_, err := io.Copy(inW, in)
		if err != nil {
			t.Logf("error copying pipe")
		}
		err = inW.CloseWithError(io.EOF)
		if err != nil {
			t.Logf("error closing pipe")
		}
		inR.CloseWithError(io.ErrClosedPipe)
	}
*/
	pipeline := NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"), OutToNull())
	//pipeline := NewPipeline(Tr("H", "O"), OutToNull())
	//pipeline := NewPipeline(Tr("H", "O"), OutToNull())
	//pipeline := NewPipeline(OutToNull())
	//go closef()
	err, code, index := pipeline.ExecAndWaitFor(NewPipeset(in, out, errout), 2 * time.Second)
	if err!=nil {
		t.Logf("Errout: %+v\n", errout.String())
		t.Logf("Stdout: %+v", out.String())
		t.Logf("Error: %+v, code: %d, index: %d\n", err, code, index)
		if code != 0 {
			t.Errorf("Error: %+v\n", err)
		}

	}
	output := out.String()
	outputErr := errout.String()
	expected := ""
	if outputErr != expected {
		t.Logf("Errout: %+v", outputErr)
		t.Logf("Stdout: %+v", output)
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}
func TestRedirectOutErrErrOut(t *testing.T) {
	var out, errout bytes.Buffer
	pipeline := NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"), OutToErr(), ErrToOut())
	err, code, index := pipeline.ExecAndWait(NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), &out, &errout))
	if err!=nil {
		t.Logf("Errout: %+v\n", errout.String())
		t.Logf("Stdout: %+v", out.String())
		t.Logf("ddError: %+v, code: %d\n", err, code)
		t.Logf("Error: %+v, code: %d, index: %d\n", err, code, index)
		if code != 0 {
			t.Errorf("Error: %+v\n", err)
		}

	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	t.Logf("Errout: %+v", errout.String())
	t.Logf("Stdout: %+v", output)
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}
