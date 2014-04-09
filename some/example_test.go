package some

import (
//	"bytes"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/wget-go/wget"
//	"strings"
)


func Example_WgetTrGzip() {
	//pipeline := someutils.NewPipeline(wget.Wget("www.golang.org"), Head(2), Tr("Go", "Gopher"), Gzip(), someutils.OutTo("test-gopher.gz"))
	
	pipeline := someutils.NewPipeline(someutils.Wrap(wget.WgetToOut("www.golang.org")), Head(2), Tr("h", "G"), Gzip(), Gunzip())
	invocation, out, errout := someutils.InvocationFromString("www.golang.org\n")
	err, invocationchan, count := invocation.PipeToPipeline(pipeline)
	if err != nil {
		fmt.Printf("error piping to pipeline: %v", err)
	}
	//err, code, index := pipeline.execandwait()
	errinvocation := someutils.Wait(invocationchan, count)
	outstring := out.String()
	if errinvocation.Err!=nil {
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

