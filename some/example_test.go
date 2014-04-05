package some

import (
	"bytes"
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/wget-go/wget"
	"strings"
)


func Example_WgetTrGzip() {
	//pipeline := someutils.NewPipeline(wget.Wget("www.golang.org"), Head(2), Tr("Go", "Gopher"), Gzip(), someutils.OutTo("test-gopher.gz"))
	
	pipeline := someutils.NewPipeline(wget.WgetToOut("www.golang.org"), Head(2), Tr("h", "G"), Gzip(), Gunzip())
	var out, errout bytes.Buffer
	err, code, index := pipeline.ExecAndWait(someutils.NewPipeset(strings.NewReader("www.golang.org\n"), &out, &errout))
	outString := out.String()
	if err!=nil {
		fmt.Printf("Errout: %+v\n", errout.String())
		fmt.Printf("Stdout: %+v", outString)
		fmt.Printf("Error: %+v, exit code: %d, index: %d\n", err, code, index)
		if code != 0 {
			fmt.Printf("Error: %+v\n", err)
		}
	}
	fmt.Println(outString)
	// Output:
	// <!DOCTYPE Gtml>
	// <Gtml>
}

