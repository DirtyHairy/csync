package local

import (
	"os"
	"path/filepath"
)

func resolveSymlink(path string) string {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		return resolved
	}

	return path
}

func statLinkTarget(path string) (os.FileInfo, error) {
	return os.Stat(resolveSymlink(path))
}
