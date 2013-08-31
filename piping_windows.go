//+build windows

package someutils

import(
	"os"
	"strings"
)
//This works for me but I have no idea if it will work for anyone else, or accross different versions of Go
func IsPipingStdin() bool {
	_, err := os.Stdin.Seek(0,1)
	isPiping := false
	if err != nil {
		if strings.Contains(err.Error(), "handle is invalid") {
			isPiping = false
		} else if strings.Contains(err.Error(), "broken pipe") {
			//seems to be true ...
			isPiping = true
		}
		//fmt.Printf("Error %v\n", err)
	} else {
		isPiping = true
	}
	return isPiping
}