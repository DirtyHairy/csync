package config

import (
	"github.com/DirtyHairy/csync/lib/storage/types"
	"github.com/blang/semver"
)

func NewMutableConfig() MutableConfig {
	return &config_v1{
		Type_: types.LOCAL,
	}
}

func NewConfig(version semver.Version) Config {
	return new(config_v1)
}
