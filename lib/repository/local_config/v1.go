package local_config

import (
	storage_config "github.com/DirtyHairy/csync/lib/storage/config"
)

type config_v1 struct {
	StorageConfig_ storage_config.Config `json:"storage"`
	Plain_         bool                  `json:"plain"`
}

func (c *config_v1) StorageConfig() storage_config.Config {
	return c.StorageConfig_
}

func (c *config_v1) Plain() bool {
	return c.Plain_
}

func (c *config_v1) SetStorageConfig(config storage_config.Config) {
	c.StorageConfig_ = config
}

func (c *config_v1) SetPlain(isPlain bool) {
	c.Plain_ = isPlain
}
