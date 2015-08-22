package storage

type DirectoryEntry interface {
	Entry

	Open() (Directory, error)
}
