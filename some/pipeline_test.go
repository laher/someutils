package some

import (
	. "github.com/laher/someutils"
	"testing"
	"time"
)

func TestPipeline1(t *testing.T) {
	pipeline := NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O")) //, RedirectTo("file.txt"))
	invocation, out, errout := InvocationFromString("Hi\nHo\nhI\nhO\n")
	err, invocationchan, count := invocation.PipeToPipeline(pipeline)
	if err != nil {
		t.Errorf("error piping to pipeline: %v", err)
	}
	errinvocation := WaitFor(invocationchan, count, 1*time.Second)
	outstring := out.String()
	if errinvocation.Err != nil {
		t.Logf("errout: %+v\n", errout.String())
		t.Logf("stdout: %+v", outstring)
		t.Logf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			t.Errorf("error: %+v\n", errinvocation.Err)
		}
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}

func TestPipeline2(t *testing.T) {
	pipeline := NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"))
	invocation, out, errout := InvocationFromString("Hi\nHo\nhI\nhO\n")
	err, invocationchan, count := invocation.PipeToPipeline(pipeline)
	if err != nil {
		t.Errorf("error piping to pipeline: %v", err)
	}
	errinvocation := Wait(invocationchan, count)
	outstring := out.String()
	if errinvocation.Err != nil {
		t.Logf("errout: %+v\n", errout.String())
		t.Logf("stdout: %+v", outstring)
		t.Logf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			t.Errorf("error: %+v\n", errinvocation.Err)
		}
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}

}

func TestRedirect1(t *testing.T) {
	pipeline := NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"), OutTo("test.txt"), Cat("test.txt"))
	invocation, out, errout := InvocationFromString("Hi\nHo\nhI\nhO\n")
	err, invocationchan, count := invocation.PipeToPipeline(pipeline)
	if err != nil {
		t.Errorf("error piping to pipeline: %v", err)
	}
	errinvocation := Wait(invocationchan, count)
	outstring := out.String()
	if errinvocation.Err != nil {
		t.Logf("errout: %+v\n", errout.String())
		t.Logf("stdout: %+v", outstring)
		t.Logf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			t.Errorf("error: %+v\n", errinvocation.Err)
		}
	}
	output := out.String()
	expected := "Oi\nOo\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}

func TestRedirectOutErr(t *testing.T) {
	pipeline := NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"), OutToErr())
	invocation, out, errout := InvocationFromString("Hi\nHo\nhI\nhO\n")
	err, invocationchan, count := invocation.PipeToPipeline(pipeline)
	if err != nil {
		t.Errorf("error piping to pipeline: %v", err)
	}
	errinvocation := Wait(invocationchan, count)
	outstring := out.String()
	if errinvocation.Err != nil {
		t.Logf("errout: %+v\n", errout.String())
		t.Logf("stdout: %+v", outstring)
		t.Logf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			t.Errorf("error: %+v\n", errinvocation.Err)
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
	pipeline := NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"), OutToNull())
	//pipeline := NewPipeline(Tr("H", "O"), OutToNull())
	//pipeline := NewPipeline(Tr("H", "O"), OutToNull())
	//pipeline := NewPipeline(OutToNull())
	//go closef()
	invocation, out, errout := InvocationFromString("Hi\nHo\nhI\nhO\n")
	err, invocationchan, count := invocation.PipeToPipeline(pipeline)
	if err != nil {
		t.Errorf("error piping to pipeline: %v", err)
	}
	errinvocation := Wait(invocationchan, count)
	outstring := out.String()
	if errinvocation.Err != nil {
		t.Logf("errout: %+v\n", errout.String())
		t.Logf("stdout: %+v", outstring)
		t.Logf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			t.Errorf("error: %+v\n", errinvocation.Err)
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
	pipeline := NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O"), OutToErr(), ErrToOut())
	invocation, out, errout := InvocationFromString("Hi\nHo\nhI\nhO\n")
	err, invocationchan, count := invocation.PipeToPipeline(pipeline)
	if err != nil {
		t.Errorf("error piping to pipeline: %v", err)
	}
	errinvocation := Wait(invocationchan, count)
	outstring := out.String()
	if errinvocation.Err != nil {
		t.Logf("errout: %+v\n", errout.String())
		t.Logf("stdout: %+v", outstring)
		t.Logf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			t.Errorf("error: %+v\n", errinvocation.Err)
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
