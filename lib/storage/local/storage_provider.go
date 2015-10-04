package local

import (
	"github.com/DirtyHairy/csync/lib/storage"
	"github.com/DirtyHairy/csync/lib/storage/local/config"

	abstract_config "github.com/DirtyHairy/csync/lib/storage/config"
)

type storageProvider struct {
	root *directory
}

func (p *storageProvider) Root() (storage.Directory, error) {
	return p.root.Entry().Open()
}

func (p *storageProvider) Marshal() abstract_config.Config {
	cfg := config.NewMutableConfig()

	cfg.SetPath(p.root.entry.realPath())

	return cfg
}
