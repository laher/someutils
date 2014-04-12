package some

import (
	"github.com/laher/someutils"
	"github.com/laher/wget-go/wget"
)

func init() {
	someutils.RegisterSimple(func() someutils.CliPipableSimple { return new(wget.Wgetter) })
}
