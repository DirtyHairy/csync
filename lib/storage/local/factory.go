package local

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DirtyHairy/csync/lib/storage"
)

func NewLocalFS(path string) (storage.Directory, error) {
	var err error

	path, err = filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(path)

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

	directory, err := entry.Open()

	if err != nil {
		return nil, err
	}

	return directory, nil
}
