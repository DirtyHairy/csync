package local

import (
	"os"

	"github.com/DirtyHairy/csync/lib/storage"
)

type fileEntry struct {
	entry
}

func (entry *fileEntry) Open() (storage.File, error) {
	file, err := os.Open(entry.realPath())

	if err != nil {
		return nil, err
	}

	return newFile(entry, file), nil
}

func (entry *fileEntry) OpenWrite() (storage.WritableFile, error) {
	file, err := os.OpenFile(entry.realPath(), os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	return newFile(entry, file), nil
}

func newFileEntry(path, prefix string, fileInfo os.FileInfo) (*fileEntry, error) {
	entry, err := newEntry(path, prefix, fileInfo)

	if err != nil {
		return nil, err
	}

	return &fileEntry{*entry}, nil
}
