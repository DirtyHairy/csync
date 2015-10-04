package local

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/DirtyHairy/csync/lib/storage"
	"github.com/DirtyHairy/csync/lib/storage/local/config"
)

func NewLocalFS(path string) (storage.StorageProvider, error) {
	var err error

	path, err = filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	fileInfo, err := statLinkTarget(path)

	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		return nil, errors.New("only directory can be used as roots for local FS adapters")
	}

	entry, err := newDirectoryEntry(fmt.Sprintf("%c", filepath.Separator), path, fileInfo)

	if err != nil {
		return nil, err
	}

	root, err := entry.Open()

	if err != nil {
		return nil, err
	}

	return &storageProvider{root.(*directory)}, nil
}

func FromConfig(cfg config.Config) (storage.StorageProvider, error) {
	return NewLocalFS(cfg.Path())
}
