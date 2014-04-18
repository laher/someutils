package some

import (
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/wget-go/wget"
	"strings"
)

func Example_WgetTr() {
	pipeline := someutils.NewPipeline(someutils.Wrap(wget.WgetToOut()), Head(2), Tr("h", "G"))
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader("www.golang.org\n"))
	errinvocation := invocation.Wait()
	outstring := out.String()
	if errinvocation.Err != nil {
		fmt.Printf("errout: %+v\n", errout.String())
		fmt.Printf("stdout: %+v", outstring)
		fmt.Printf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			fmt.Printf("error: %+v\n", errinvocation.Err)
		}
	}
	fmt.Println(outstring)
	// Output:
	// <!DOCTYPE Gtml>
	// <Gtml>
}

/*
for another test ...
	//pipeline := someutils.NewPipeline(wget.Wget("www.golang.org"), Head(2), Tr("G", "O"), Gzip(), someutils.OutTo("test-gopher.gz"))
*/

func Example_GzipGunzip() {
	pipeline := someutils.NewPipeline(Gzip(), Gunzip())
	input := "hiya"
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader(input))
	errinvocation := invocation.Wait()
	outstring := out.String()
	if errinvocation.Err != nil {
		fmt.Printf("errout: %+v\n", errout.String())
		fmt.Printf("stdout: %+v", outstring)
		fmt.Printf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			fmt.Printf("error: %+v\n", errinvocation.Err)
		}
	}
	fmt.Println(outstring)
	// Output:
	// hiya
}

func Example_WgetHeadTr() {
	pipeline := someutils.NewPipeline(someutils.Wrap(wget.WgetToOut()), Head(2), Tr("h", "G")) //, Gzip()) //, Gunzip()) //, Tr("h", "G"))
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader("www.golang.org\n"))
	errinvocation := invocation.Wait()
	outstring := out.String()
	if errinvocation.Err != nil {
		fmt.Printf("errout: %+v\n", errout.String())
		fmt.Printf("stdout: %+v\n", outstring)
		fmt.Printf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			fmt.Printf("error: %+v\n", errinvocation.Err)
		}
	}
	fmt.Println(outstring)
	// Output:
	// <!DOCTYPE Gtml>
	// <Gtml>
}

func Example_TrNoFanout() {
	pipeline := someutils.NewPipeline(Tr("w", "v"))
	input := "www.google.com\nwww.bbc.co.uk\nwww.golang.org\n"
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader(input))
	errinvocation := invocation.Wait()
	outstring := out.String()
	if errinvocation.Err != nil {
		fmt.Printf("errout: %+v\n", errout.String())
		fmt.Printf("stdout: %+v", outstring)
		fmt.Printf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			fmt.Printf("error: %+v\n", errinvocation.Err)
		}
	}
	fmt.Println(outstring)
	// Output:
	// vvv.google.com
	// vvv.bbc.co.uk
	// vvv.golang.org
}

func Example_TrFanout() {
	pipeline := someutils.NewPipeline(someutils.FanoutByLine(Tr("w", "v")))
	input := "www.google.com\nwww.bbc.co.uk\nwww.golang.org\n"
	invocation, out, errout := pipeline.InvokeReader(strings.NewReader(input))
	errinvocation := invocation.Wait()
	outstring := out.String()
	if errinvocation.Err != nil {
		fmt.Printf("errout: %+v\n", errout.String())
		fmt.Printf("stdout: %+v", outstring)
		fmt.Printf("error: %+v, exit code: %d\n", errinvocation.Err, errinvocation.ExitCode)
		if *errinvocation.ExitCode != 0 {
			fmt.Printf("error: %+v\n", errinvocation.Err)
		}
	}
	fmt.Println(outstring)
	// Output:
	// vvv.google.com
	// vvv.bbc.co.uk
	// vvv.golang.org
}

