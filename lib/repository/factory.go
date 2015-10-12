package repository

import (
	"errors"

	"github.com/DirtyHairy/csync/lib/repository/local_config"

	storage_factory "github.com/DirtyHairy/csync/lib/storage/factory"
)

func CreateFromConfig(cfg local_config.Config) (repo Repository, err error) {
	if !cfg.Plain() {
		err = errors.New("non-plain repositories are not yet supported")
		return
	}

	underlyingStorage, err := storage_factory.FromConfig(cfg.StorageConfig())

	if err != nil {
		return
	}

	repo = NewRepository(cfg.Plain(), underlyingStorage)
	return
}
