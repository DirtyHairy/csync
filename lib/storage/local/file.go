package local

import (
	"os"

	"github.com/DirtyHairy/csync/lib/storage"
)

type file struct {
	entry  *fileEntry
	file   *os.File
	closed bool
}

func (f *file) Entry() storage.FileEntry {
	return f.entry
}

func (f *file) Read(buf []byte) (int, error) {
	return f.file.Read(buf)
}

func (f *file) Write(buffer []byte) (int, error) {
	return f.file.Write(buffer)
}

func (f *file) Close() error {
	if f.closed {
		return nil
	}

	err := f.file.Close()

	if err != nil {
		return err
	}

	f.closed = true

	return nil
}

func newFile(entry *fileEntry, f *os.File) *file {
	return &file{
		entry: entry,
		file:  f,
	}
}
