//TODO test with freebsd,netbsd

// +build darwin freebsd netbsd

package someutils

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA