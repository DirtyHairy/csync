package local

import (
	"errors"
	"os"
	"path/filepath"
	"time"
	"unicode/utf8"
)

type entry struct {
	path     string
	prefix   string
	fileInfo os.FileInfo
}

func (e *entry) Name() string {
	return e.fileInfo.Name()
}

func (e *entry) Path() string {
	return filepath.ToSlash(e.path)
}

func (e *entry) Mtime() time.Time {
	return e.fileInfo.ModTime()
}

func (e *entry) realPath() string {
	return filepath.Join(e.prefix, e.path)
}

func newEntry(path, prefix string, fileInfo os.FileInfo) (*entry, error) {
	if !utf8.ValidString(path) || !utf8.ValidString(prefix) || !utf8.ValidString(fileInfo.Name()) {
		return nil, errors.New("pathname is not valid UTF8")
	}

	entry := entry{
		path:     path,
		prefix:   prefix,
		fileInfo: fileInfo,
	}

	return &entry, nil
}
