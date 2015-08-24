// +build linux darwin windows netbsd openbsd freebsd dragonfly

package local

import (
	"os"
	"syscall"
	"time"
)

// This implementation should work for all systems that implement the utimes syscall
func (e *entry) SetMtime(mtime time.Time) error {
	tvMtime := syscall.NsecToTimeval(mtime.UnixNano())
	tvAtime := syscall.NsecToTimeval(time.Now().UnixNano())

	resolvedPath := resolveSymlink(e.realPath())

	err := syscall.Utimes(resolvedPath, []syscall.Timeval{tvAtime, tvMtime})

	if err != nil {
		return err
	}

	e.fileInfo, err = os.Stat(resolvedPath)

	if err != nil {
		return err
	}

	return nil
}
