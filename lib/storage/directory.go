package storage

type Directory interface {
	Entry() DirectoryEntry

	Stat(path string) (Entry, error)

	NextEntry() (Entry, error)

	Rewind() error

	CreateFile(path string) (WritableFile, error)

	Mkdir(path string) (Directory, error)

	Close() error
}
