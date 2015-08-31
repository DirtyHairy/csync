package storage

import (
	"io"
)

type Directory interface {
	io.Closer

	Entry() DirectoryEntry

	Stat(path string) (Entry, error)

	NextEntry() (Entry, error)

	Rewind() error

	CreateFile(path string) (WritableFile, error)

	Mkdir(path string) (Directory, error)
}
