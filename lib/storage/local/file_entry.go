package local

import (
	"errors"
	"os"
	"path/filepath"

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

func (entry *fileEntry) Rename(newname string) (storage.Entry, error) {
	if len(filepath.SplitList(filepath.FromSlash(newname))) > 1 {
		return nil, errors.New("Rename does not accept paths")
	}

	pathDir, _ := filepath.Split(entry.Path())
	newPath := filepath.Join(pathDir, newname)
	newRealPath := filepath.Join(entry.prefix, newPath)

	if err := os.Rename(entry.realPath(), newRealPath); err != nil {
		return nil, err
	}

	fi, err := statLinkTarget(newRealPath)

	if err != nil {
		return nil, err
	}

	return newFileEntry(newPath, entry.prefix, fi)
}

func newFileEntry(path, prefix string, fileInfo os.FileInfo) (*fileEntry, error) {
	entry, err := newEntry(path, prefix, fileInfo)

	if err != nil {
		return nil, err
	}

	return &fileEntry{*entry}, nil
}
