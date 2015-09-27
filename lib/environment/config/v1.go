package config

import (
	"github.com/blang/semver"
)

type config_v1 struct {
	Version semver.Version `json:"csyncVersion"`
}

func (c *config_v1) CsyncVersion() semver.Version {
	return c.Version
}
