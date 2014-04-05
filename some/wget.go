package some

import (
	"github.com/laher/someutils"
	"github.com/laher/wget-go/wget"
)

func init() {
	someutils.RegisterPipable(func() someutils.NamedPipable { return new (wget.Wgetter) })
}

