package local

import (
	"github.com/DirtyHairy/csync/lib/storage"
)

type storageProvider struct {
	root *directory
}

func (p *storageProvider) Root() (storage.Directory, error) {
	return p.root.Entry().Open()
}
