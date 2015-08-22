package local

import (
	"os"

	"github.com/DirtyHairy/csync/lib/storage"
)

type directoryEntry struct {
	entry
}

func (d *directoryEntry) Open() (storage.Directory, error) {
	return newDirectory(d), nil
}

func newDirectoryEntry(path, prefix string, fileInfo os.FileInfo) (*directoryEntry, error) {
	entry, err := newEntry(path, prefix, fileInfo)

	if err != nil {
		return nil, err
	}

	return &directoryEntry{*entry}, nil
}
