package sync

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/DirtyHairy/csync/lib/storage"
)

type Unidirectional struct {
	config Config
	rng    *rand.Rand
}

func (sync *Unidirectional) tempFileNameFor(entry storage.Entry) string {
	hash := sha1.New()

	_, _ = hash.Write([]byte(entry.Path()))
	_, _ = hash.Write([]byte(time.Now().String()))

	randomBytes := make([]byte, 32)
	for i := 0; i < 32; i++ {
		randomBytes[i] = byte(sync.rng.Int())
	}

	_, _ = hash.Write(randomBytes)

	buffer := hash.Sum(make([]byte, 0, hash.Size()))

	return "csync_temp_" + hex.EncodeToString(buffer)
}

func (sync *Unidirectional) syncFailedError(originalError error) error {
	return errors.New(fmt.Sprintf("sync failed: %v", originalError))
}

func (sync *Unidirectional) syncFile(dirEntry storage.DirectoryEntry, entry storage.FileEntry) error {
	to := sync.config.To
	path := entry.Path()
	mustSync := false
	targetEntry, err := to.Stat(path)

	if err != nil {
		return err
	}

	if targetEntry == nil {
		mustSync = true
	} else if _, ok := targetEntry.(storage.FileEntry); !ok {
		mustSync = true
	} else if targetEntry.Mtime().Before(entry.Mtime()) {
		mustSync = true
	}

	if !mustSync {
		return nil
	}

	targetDirEntry, err := sync.config.To.Stat(dirEntry.Path())

	if err != nil {
		return err
	}

	targetDir, err := targetDirEntry.(storage.DirectoryEntry).Open()

	if err != nil {
		return err
	}

	defer func() {
		_ = targetDir.Close()
	}()

	tempFile, err := targetDir.CreateFile(sync.tempFileNameFor(entry))

	if err != nil {
		return err
	}

	defer func() {
		_ = tempFile.Close()
		_ = tempFile.Entry().Remove()
	}()

	soureFile, err := entry.Open()

	if err != nil {
		return err
	}

	defer func() {
		_ = soureFile.Close()
	}()

	fmt.Printf("copying %s to %s\n", entry.Path(), tempFile.Entry().Path())

	if _, err := io.Copy(tempFile, soureFile); err != nil {
		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	if err := soureFile.Close(); err != nil {
		return err
	}

	if err := tempFile.Entry().SetMtime(entry.Mtime()); err != nil {
		return err
	}

	if targetEntry != nil {
		fmt.Printf("removing %s\n", targetEntry.Path())

		if err := targetEntry.Remove(); err != nil {
			return err
		}
	}

	fmt.Printf("replacing %s with %s\n", entry.Path(), tempFile.Entry().Path())

	if _, err := tempFile.Entry().Rename(entry.Name()); err != nil {
		return err
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
				err := sync.syncFile(dir.Entry(), entry)

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
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
