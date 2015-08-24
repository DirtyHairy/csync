package local

import (
	"os"
)

func resolveSymlink(path string) string {
	var err error

	for err == nil {
		var newpath string

		newpath, err = os.Readlink(path)

		if err == nil {
			path = newpath
		}
	}

	return path
}

func statLinkTarget(path string) (os.FileInfo, error) {
	return os.Stat(resolveSymlink(path))
}
