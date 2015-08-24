package sync

import (
	"errors"
	"fmt"
	"io"

	"github.com/DirtyHairy/csync/lib/storage"
)

type Unidirectional struct {
	config Config
}

func (sync *Unidirectional) syncFailedError(originalError error) error {
	return errors.New(fmt.Sprintf("sync failed: %v", originalError))
}

func (sync *Unidirectional) syncFile(entry storage.FileEntry) error {
	to := sync.config.To
	path := entry.Path()

	var (
		targetEntry     storage.Entry
		targetFileEntry storage.FileEntry
		targetFile      storage.WritableFile
		err             error
		ok              bool
	)

	targetEntry, err = to.Stat(path)

	if err != nil {
		return err
	}

	if targetFileEntry, ok = targetEntry.(storage.FileEntry); !ok && targetEntry != nil {
		fmt.Printf("removing %s in target repo\n", path)

		if err := targetEntry.Remove(); err != nil {
			return err
		}

		targetFileEntry = nil
	}

	if targetFileEntry == nil {
		targetFile, err = to.CreateFile(path)
	} else {
		if targetFileEntry.Mtime().Before(entry.Mtime()) {
			targetFile, err = targetFileEntry.OpenWrite()
		}
	}

	if err != nil {
		return err
	}

	if targetFile != nil {
		fmt.Printf("copying %s\n", path)

		defer func() {
			_ = targetFile.Close()
		}()

		sourceFile, err := entry.Open()

		if err != nil {
			return err
		}

		defer func() {
			_ = sourceFile.Close()
		}()

		if _, err = io.Copy(targetFile, sourceFile); err != nil {
			return err
		}

		if err := targetFile.Entry().SetMtime(sourceFile.Entry().Mtime()); err != nil {
			return err
		}
	}

	return nil
}

func (sync *Unidirectional) syncDir(entry storage.DirectoryEntry) error {
	to := sync.config.To
	path := entry.Path()

	var (
		targetEntry    storage.Entry
		targetDirEntry storage.DirectoryEntry
		targetDir      storage.Directory
		err            error
		ok             bool
	)

	targetEntry, err = to.Stat(path)

	if err != nil {
		return err
	}

	if targetDirEntry, ok = targetEntry.(storage.DirectoryEntry); !ok && targetEntry != nil {
		fmt.Printf("removing %s in target repo\n", path)

		if err := targetEntry.Remove(); err != nil {
			return err
		}

		targetDirEntry = nil
	}

	if targetDirEntry != nil {
		return nil
	}

	fmt.Printf("creating directory %s\n", path)

	targetDir, err = to.Mkdir(path)

	if err != nil {
		return err
	}

	defer func() {
		targetDir.Close()
	}()

	if err := targetDir.Entry().SetMtime(entry.Mtime()); err != nil {
		return err
	}

	return nil
}

func (sync *Unidirectional) Execute() error {
	directoryStack := make([]storage.Directory, 0, 10)
	directoryStack = append(directoryStack, sync.config.From)

	for len(directoryStack) > 0 {
		dir := directoryStack[0]
		directoryStack = directoryStack[1:]

		var (
			err   error
			entry storage.Entry
		)

		for entry, err = dir.NextEntry(); entry != nil; entry, err = dir.NextEntry() {
			if err != nil {
				return sync.syncFailedError(err)
			}

			switch entry := entry.(type) {

			case storage.DirectoryEntry:
				err := sync.syncDir(entry)

				if err != nil {
					return sync.syncFailedError(err)
				}

				dir, err := entry.Open()

				if err != nil {
					return sync.syncFailedError(err)
				}

				directoryStack = append(directoryStack, dir)

			case storage.FileEntry:
				err := sync.syncFile(entry)

				if err != nil {
					return sync.syncFailedError(err)
				}
			}
		}

		if err != nil {
			return sync.syncFailedError(err)
		}
	}

	return nil
}

func NewUnidirectionalSync(config Config) *Unidirectional {
	return &Unidirectional{
		config: config,
	}
}
