package config

import (
	"github.com/blang/semver"
)

func NewConfig(csyncVersion semver.Version) Config {
	return &config_v1{
		Version: csyncVersion,
	}
}

func emptyConfigForVersion(CsyncVersion semver.Version) Config {
	return new(config_v1)
}
