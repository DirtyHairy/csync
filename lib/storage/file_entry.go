package storage

type FileEntry interface {
	Entry

	Open() (File, error)
	OpenWrite() (WritableFile, error)
}
