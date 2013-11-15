//+build !windows

package someutils

import (
	"os"
)

func makeLink(target, linkName string, options LnOptions) error {
	if options.IsSymbolic {
		return os.Symlink(target, linkName)
	} else {
		return os.Link(target, linkName)
	}
}
