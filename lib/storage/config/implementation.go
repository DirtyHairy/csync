package config

import (
	"encoding/json"
	"errors"

	config_local "github.com/DirtyHairy/csync/lib/storage/local/config"
	"github.com/DirtyHairy/csync/lib/storage/types"
	"github.com/blang/semver"
)

// Wraps the config specific to a particular storage provider
type unspecificConfig struct {
	Config
	forVersion semver.Version
}

// Try to identify the storage provider by type and unmarshal the wrapper
// accordingly. The result is suitable for storage/factory .
func (m *unspecificConfig) UnmarshalJSON(jsonData []byte) (err error) {
	var protoConfig struct {
		Type string `json:"type"`
	}

	err = json.Unmarshal(jsonData, &protoConfig)
	if err != nil {
		return
	}

	switch protoConfig.Type {
	case types.LOCAL:
		m.Config = config_local.NewConfig(m.forVersion)

	default:
		err = errors.New("invalid type")
		return
	}

	err = json.Unmarshal(jsonData, m.Config)

	return
}

// Unwrap config in case we are dealing with a unspecificConfig. This cludge is
// necessary as storage/factory has no knowledge of the inner structure of this
// type. Moving the factory to this package is no option as there will be a
// cyclic dependency of the type
//
// storage/config -> storage (via return type) -> storage/config (via StorageProvider.Marshal)
func Unwrap(in Config) Config {
	switch in := in.(type) {
	case *unspecificConfig:
		return in.Config

	default:
		return in
	}
}
