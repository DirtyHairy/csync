package environment

import (
	"github.com/blang/semver"
)

type VersionProvider interface {
	Version() semver.Version
}

type VersionReceiver interface {
	SetVersion(semver.Version)
}

type Environment interface {
	VersionProvider

	Save() error
}

type MutableEnvironment interface {
	Environment

	VersionReceiver

	Load() error
}
