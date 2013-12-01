//+build !freebsd,!netbsd,!openbsd,!plan9

package someutils

import (
	"github.com/laher/scp-go/scp"
)

func init() {
	Register(Util{
		"scp",
		scp.Scp})
}
