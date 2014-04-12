//+build !freebsd,!netbsd,!openbsd,!plan9

package some

import (
	"github.com/laher/scp-go/scp"
	"github.com/laher/someutils"
)

func init() {
	someutils.RegisterSimple(func() someutils.CliPipableSimple { return new(scp.SecureCopier) })
}

