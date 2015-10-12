package repository

import (
	"github.com/DirtyHairy/csync/lib/repository/local_config"
	"github.com/DirtyHairy/csync/lib/storage"
)

type Repository interface {
	Root() storage.Directory

	Initialize() error

	Marshal() local_config.Config
}
