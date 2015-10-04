// This needs to be a separate package in order to avoid cyclic dependencies.

package factory

import (
	"errors"
	"fmt"

	"github.com/DirtyHairy/csync/lib/storage"
	"github.com/DirtyHairy/csync/lib/storage/local"
	"github.com/DirtyHairy/csync/lib/storage/types"

	config_unspecific "github.com/DirtyHairy/csync/lib/storage/config"
	config_local "github.com/DirtyHairy/csync/lib/storage/local/config"
)

// Create a storage provider from the provided config.
func FromConfig(cfg config_unspecific.Config) (provider storage.StorageProvider, err error) {
	// Make sure that we are dealing with a specific config
	cfg = config_unspecific.Unwrap(cfg)

	switch cfg.Type() {
	case types.LOCAL:
		provider, err = local.FromConfig(cfg.(config_local.Config))

	default:
		err = errors.New(fmt.Sprintf(`cannot happen: invalid storage type "%s"`, cfg.Type()))
	}

	return
}
