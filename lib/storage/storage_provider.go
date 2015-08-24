package storage

type StorageProvider interface {
	Root() (Directory, error)
}
