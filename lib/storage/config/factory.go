package config

import (
	"github.com/blang/semver"
)

func NewConfig(version semver.Version) Config {
	return &unspecificConfig{
		forVersion: version,
	}
}
