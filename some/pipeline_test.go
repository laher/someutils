package some

import (
	. "github.com/laher/someutils"
	"strings"
	"testing"
	"time"
)

func TestPipeline1(t *testing.T) {
	pipeline := NewPipeline(Tr("H", "O"), Tr("I", "J"), Grep("O")) //, RedirectTo("file.txt"))
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader("Hi\nHo\nhI\nhO\n"))
	errinvocation := invocation.WaitUpTo(1 * time.Second)
	outstring := out.String()
	if errinvocation == nil || errinvocation.Err != nil {
		t.Logf("errout: %+v\n", errout.String())
		t.Logf("stdout: %+v", outstring)
		if errinvocation == nil {
			t.Errorf("errinvocation is nil!!!\n")

		} else {
			t.Logf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
			if *errinvocation.ExitCode != 0 {
				t.Errorf("error: %+v\n", errinvocation.Err)
			}
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
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader("Hi\nHo\nhI\nhO\n"))
	errinvocation := invocation.Wait()
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
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader("Hi\nHo\nhI\nhO\n"))
	errinvocation := invocation.Wait()
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
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader("Hi\nHo\nhI\nhO\n"))
	errinvocation := invocation.Wait()
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
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader("Hi\nHo\nhI\nhO\n"))
	errinvocation := invocation.Wait()
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
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader("Hi\nHo\nhI\nhO\n"))
	errinvocation := invocation.Wait()
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
		t.Logf("Errout: %+v", errout.String())
		t.Logf("Stdout: %+v", output)
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
}
