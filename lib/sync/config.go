package sync

import (
	"github.com/DirtyHairy/csync/lib/storage"
)

type Config struct {
	From storage.StorageProvider
	To   storage.StorageProvider
}
