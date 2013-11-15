package someutils

import (
	"github.com/laher/wget-go/wget"
)

func init() {
	Register(Util{
		"wget",
		wget.Wget})
}
