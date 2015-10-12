package config

import (
	"github.com/blang/semver"
)

func NewMutableConfig() MutableConfig {
	// the map will allocated on-the-fly, so we can use the zero type here
	return new(config_v1)
}

func NewConfig(csyncVersion semver.Version) Config {
	return &config_v1{
		// UnmarshalJSON will allocate the map on the fly
		forVersion: csyncVersion,
	}
}
