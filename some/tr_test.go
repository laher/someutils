package some

import (
	"bytes"
	"fmt"
	"github.com/laher/someutils"
	"io"
	"strings"
	"testing"
	"time"
)

func TestTrCli(t *testing.T) {
	//var out bytes.Buffer
	//var errout bytes.Buffer
	//pipes := someutils.NewPipes(strings.NewReader("HI"), &out, &errout)
	err := TrCli([]string{"tr", "I", "O"})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	//println(out.String())
}

func TestFluentTr(t *testing.T) {
	var out bytes.Buffer
	var errout bytes.Buffer
	pipes := someutils.NewPipes(strings.NewReader("HI"), &out, &errout)
	err := Tr("I", "O").Exec(pipes)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	println(out.String())
}

func Test2pipes(t *testing.T) {
	var out bytes.Buffer
	var errout bytes.Buffer
	r, w := io.Pipe()
	in := strings.NewReader("Hi\nHo\nhI\nhO\n")
	pipes1 := someutils.NewPipes(in, w, &errout)
	pipes2 := someutils.NewPipes(r, &out, &errout)
	tr1 := Tr("H", "O")
	tr2 := Tr("I", "J")
	go tr1.Exec(pipes1)
	go tr2.Exec(pipes2)
	time.Sleep(1 * time.Second)
	output := out.String()
	expected := "Oi\nOo\nhJ\nhO\n"
	if output != expected {
		t.Error("Expected ", expected, ", got ", output)
	}
	fmt.Printf("Errout: %+v\n", errout.String())
}

func TestPipeline(t *testing.T) {
	var out bytes.Buffer
	var errout bytes.Buffer
	in := strings.NewReader("Hi\nHo\nhI\nhO\n")
	pipes := someutils.NewPipes(in, &out, &errout)
	e := someutils.Pipeline(pipes, Tr("H", "O"), Tr("I", "J"))
	errs := someutils.CollectErrors(e, 2)
	output := out.String()
	expected := "Oi\nOo\nhJ\nhO\n"
	if output != expected {
		t.Error("Expected\n ", expected, ", Got:\n ", output)
	}
	fmt.Printf("Errors: %d, %+v\n", len(errs), errs)
	fmt.Printf("Errout: %+v\n", errout.String())
}
