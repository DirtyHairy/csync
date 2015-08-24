// +build linux darwin windows netbsd openbsd freebsd dragonfly

package local

import (
	"os"
	"syscall"
	"time"
)

func (e *entry) SetMtime(mtime time.Time) error {
	tvMtime := syscall.NsecToTimeval(mtime.UnixNano())
	tvAtime := syscall.NsecToTimeval(time.Now().UnixNano())

	err := syscall.Utimes(e.realPath(), []syscall.Timeval{tvAtime, tvMtime})

	if err != nil {
		return err
	}

	e.fileInfo, err = os.Stat(e.realPath())

	if err != nil {
		return err
	}

	return nil
}
