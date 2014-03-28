package some

import (
	"bytes"
	"fmt"
	"github.com/laher/someutils"
	"strings"
	"testing"
)

func TestPipeline1(t *testing.T) {
	var out bytes.Buffer
	var errout bytes.Buffer
	var errin bytes.Buffer
	in := strings.NewReader("Hi\nHo\nhI\nhO\n")
	p := someutils.Pipeline{in, &out, &errout, &errin}
	ok, errs := p.PipeAndWait(1, Tr("H", "O"), Tr("I", "J"), Grep("O"))//, RedirectTo("file.txt"))
	if !ok {
		fmt.Printf("Errors: %d, %+v\n", len(errs), errs)
		fmt.Printf("Errout: %+v\n", errout.String())
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}

func TestPipeline2(t *testing.T) {
	p, out, errout := someutils.NewPipelineFromString("Hi\nHo\nhI\nhO\n")
	ok, errs := p.PipeAndWait(1, Tr("H", "O"), Tr("I", "J"), Grep("O"))
	if !ok {
		t.Logf("Errors: %d, %+v\n", len(errs), errs)
		t.Logf("Errout: %+v\n", errout.String())
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}

}


func TestRedirect1(t *testing.T) {
	p, out, errout := someutils.NewPipelineFromString("Hi\nHo\nhI\nhO\n")
	ok, errs := p.PipeAndWait(1, Tr("H", "O"), Tr("I", "J"), Grep("O"), RedirectTo("file.txt"), Cat("file.txt"))
	if !ok {
		t.Logf("Errors: %d, %+v\n", len(errs), errs)
		t.Logf("Errout: %+v\n", errout.String())
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}

func TestRedirectOutErr(t *testing.T) {
	p, out, errout := someutils.NewPipelineFromString("Hi\nHo\nhI\nhO\n")
	ok, errs := p.PipeAndWait(1, Tr("H", "O"), Tr("I", "J"), Grep("O"), RedirectOutToErr())
	if !ok {
		t.Logf("Errors: %d, %+v\n", len(errs), errs)
		t.Logf("Errout: %+v\n", errout.String())
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


func TestRedirectOutErrErrOut(t *testing.T) {
	p, out, errout := someutils.NewPipelineFromString("Hi\nHo\nhI\nhO\n")
	ok, errs := p.PipeAndWait(6, Tr("H", "O"), Tr("I", "J"), Grep("O"), RedirectOutToErr(), RedirectErrToOut())
	if !ok {
		t.Logf("Errors: %d, %+v\n", len(errs), errs)
		t.Logf("Errout: %+v\n", errout.String())
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
		t.Logf("Errout: %+v", errout.String())
		t.Logf("Stdout: %+v", output)
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}


