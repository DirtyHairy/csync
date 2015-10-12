package config

import (
	"github.com/DirtyHairy/csync/lib/repository/local_config"
)

type Config interface {
	RepositoryConfigs() map[string]local_config.Config
}

type MutableConfig interface {
	Config

	AddConfig(key string, cfg local_config.Config)
}
