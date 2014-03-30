package some

import (
	"bytes"
	"github.com/laher/someutils"
	"strings"
	"testing"
	"time"
)

func TestPipeline1(t *testing.T) {
	var out, errout bytes.Buffer
	pipeline := someutils.NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O")) //, RedirectTo("file.txt"))
	ok, errs := pipeline.PipeAndWait(someutils.NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), &out, &errout))
	if !ok {
		t.Logf("Errout: %+v\n", errout.String())
		t.Logf("Stdout: %+v", out.String())
		t.Errorf("Errors: %d, %+v\n", len(errs), errs)
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}

func TestPipeline2(t *testing.T) {
	var out, errout bytes.Buffer
	p := someutils.NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"))
	ok, errs := p.PipeAndWait(someutils.NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), &out, &errout))
	if !ok {
		t.Logf("Errout: %+v\n", errout.String())
		t.Logf("Stdout: %+v", out.String())
		t.Errorf("Errors: %d, %+v\n", len(errs), errs)
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}

}

func TestRedirect1(t *testing.T) {
	var out, errout bytes.Buffer
	p := someutils.NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"), OutTo("test.txt"), Cat("test.txt"))
	ok, errs := p.PipeAndWait(someutils.NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), &out, &errout))
	if !ok {
		t.Logf("Errout: %+v\n", errout.String())
		t.Logf("Stdout: %+v", out.String())
		t.Errorf("Errors: %d, %+v\n", len(errs), errs)
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}

func TestRedirectOutErr(t *testing.T) {
	var out, errout bytes.Buffer
	p := someutils.NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"), OutToErr())
	ok, errs := p.PipeAndWait(someutils.NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), &out, &errout))
	if !ok {
		t.Logf("Errout: %+v\n", errout.String())
		t.Logf("Stdout: %+v", out.String())
		t.Errorf("Errors: %d, %+v\n", len(errs), errs)
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
	var out, errout bytes.Buffer
	p := someutils.NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"), OutToNull())
	ok, errs := p.PipeAndWaitFor(someutils.NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), &out, &errout), 1 * time.Second)
	if !ok {
		t.Logf("Errout: %+v\n", errout.String())
		t.Logf("Stdout: %+v", out.String())
		t.Errorf("Errors: %d, %+v\n", len(errs), errs)
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
	p := someutils.NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"), OutToErr(), ErrToOut())
	ok, errs := p.PipeAndWait(someutils.NewPipeset(strings.NewReader("Hi\nHo\nhI\nhO\n"), &out, &errout))
	if !ok {
		t.Logf("Errout: %+v\n", errout.String())
		t.Logf("Stdout: %+v", out.String())
		t.Errorf("Errors: %d, %+v\n", len(errs), errs)
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	t.Logf("Errout: %+v", errout.String())
	t.Logf("Stdout: %+v", output)
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}
