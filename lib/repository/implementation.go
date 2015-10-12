package repository

import (
	"errors"

	"github.com/DirtyHairy/csync/lib/repository/local_config"
	"github.com/DirtyHairy/csync/lib/storage"
)

type repository struct {
	underlyingStorage storage.StorageProvider

	plain bool

	root storage.Directory
}

func (r *repository) Root() storage.Directory {
	return r.root
}

func (r *repository) Initialize() (err error) {
	if !r.plain {
		err = errors.New("non-plain repositories are not implemented yet")
		return
	}

	r.root, err = r.underlyingStorage.Root()

	return
}

func (r *repository) Marshal() local_config.Config {
	cfg := local_config.NewMutableConfig()

	cfg.SetPlain(r.plain)
	cfg.SetStorageConfig(r.underlyingStorage.Marshal())

	return cfg
}

func NewRepository(plain bool, storage storage.StorageProvider) Repository {
	return &repository{
		plain:             plain,
		underlyingStorage: storage,
	}
}
