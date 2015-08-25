package local

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/DirtyHairy/csync/lib/storage"
)

type directoryEntry struct {
	entry
}

func (d *directoryEntry) Open() (storage.Directory, error) {
	return newDirectory(d), nil
}

func (d *directoryEntry) Rename(newname string) (storage.Entry, error) {
	if len(filepath.SplitList(filepath.FromSlash(newname))) > 1 {
		return nil, errors.New("Rename does not accept paths")
	}

	pathDir, _ := filepath.Split(d.Path())
	newPath := filepath.Join(pathDir, newname)
	newRealPath := filepath.Join(d.prefix, newname)

	if err := os.Rename(d.realPath(), newRealPath); err != nil {
		return nil, err
	}

	fi, err := statLinkTarget(newRealPath)

	if err != nil {
		return nil, err
	}

	return newDirectoryEntry(newPath, d.prefix, fi)
}

func newDirectoryEntry(path, prefix string, fileInfo os.FileInfo) (*directoryEntry, error) {
	entry, err := newEntry(path, prefix, fileInfo)

	if err != nil {
		return nil, err
	}

	return &directoryEntry{*entry}, nil
}
