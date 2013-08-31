// +build linux

package someutils

import "syscall"

const ioctlReadTermios = syscall.TCGETS