package storage

import (
	"github.com/DirtyHairy/csync/lib/storage/config"
)

type StorageProvider interface {
	Root() (Directory, error)

	Marshal() config.Config
}
