package local_config

import (
	storage_config "github.com/DirtyHairy/csync/lib/storage/config"

	"github.com/blang/semver"
)

func NewMutableConfig() MutableConfig {
	return new(config_v1)
}

func NewConfig(csyncVersion semver.Version) Config {
	return &config_v1{
		StorageConfig_: storage_config.NewConfig(csyncVersion),
	}
}
