package config

import (
	"github.com/blang/semver"
)

type Config interface {
	CsyncVersion() semver.Version
}

type MutableConfig interface {
	Config
}
