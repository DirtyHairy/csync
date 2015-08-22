package local

import (
	"os"

	"github.com/DirtyHairy/csync/lib/storage"
)

type file struct {
	entry *fileEntry
	file  *os.File
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
	return f.file.Close()
}

func newFile(entry *fileEntry, f *os.File) *file {
	return &file{
		entry: entry,
		file:  f,
	}
}
