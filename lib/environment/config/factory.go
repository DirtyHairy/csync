package config

import (
	"github.com/blang/semver"
)

func NewMutableConfig(csyncVersion semver.Version) MutableConfig {
	return &config_v1{
		Version: csyncVersion,
	}
}

func newConfig(CsyncVersion semver.Version) Config {
	return new(config_v1)
}
