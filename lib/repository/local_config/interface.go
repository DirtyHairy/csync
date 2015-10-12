package local_config

import (
	storage_config "github.com/DirtyHairy/csync/lib/storage/config"
)

type Config interface {
	StorageConfig() storage_config.Config
	Plain() bool
}

type MutableConfig interface {
	Config

	SetStorageConfig(storage_config.Config)
	SetPlain(bool)
}
