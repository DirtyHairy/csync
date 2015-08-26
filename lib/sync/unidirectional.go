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

	from storage.Directory
	to   storage.Directory

	state int

	directoryStack   []storage.Directory
	currentDirectory storage.Directory
	currentlySyncing storage.Entry

	lastError error

	events chan EventInterface
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

func (sync *Unidirectional) setLastErrorIfApplicable(err error) {
	if sync.lastError == nil {
		sync.lastError = err
	}
}

func (sync *Unidirectional) syncFile(dirEntry storage.DirectoryEntry, entry storage.FileEntry) error {
	path := entry.Path()
	mustSync := false
	targetEntry, err := sync.to.Stat(path)

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

	sync.events <- &EventStartSyncing{entry: entry}

	targetDirEntry, err := sync.to.Stat(dirEntry.Path())

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
		entry, err := sync.to.Stat(tempFile.Entry().Path())

		if err != nil {
			sync.setLastErrorIfApplicable(err)
			return
		}

		if entry != nil {
			sync.setLastErrorIfApplicable(tempFile.Close())
			sync.setLastErrorIfApplicable(tempFile.Entry().Remove())
		}
	}()

	soureFile, err := entry.Open()

	if err != nil {
		return err
	}

	defer func() {
		sync.setLastErrorIfApplicable(soureFile.Close())
	}()

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
		if err := targetEntry.Remove(); err != nil {
			return err
		}
	}

	if _, err := tempFile.Entry().Rename(entry.Name()); err != nil {
		return err
	}

	return nil
}

func (sync *Unidirectional) syncDir(entry storage.DirectoryEntry) error {
	path := entry.Path()

	var (
		targetEntry    storage.Entry
		targetDirEntry storage.DirectoryEntry
		targetDir      storage.Directory
		err            error
		ok             bool
	)

	targetEntry, err = sync.to.Stat(path)

	if err != nil {
		return err
	}

	if targetDirEntry, ok = targetEntry.(storage.DirectoryEntry); !ok && targetEntry != nil {
		if err := targetEntry.Remove(); err != nil {
			return err
		}

		targetDirEntry = nil
	}

	if targetDirEntry != nil {
		return nil
	}

	sync.events <- &EventStartSyncing{entry: entry}

	targetDir, err = sync.to.Mkdir(path)

	if err != nil {
		return err
	}

	defer func() {
		sync.setLastErrorIfApplicable(targetDir.Close())
	}()

	return nil
}

func (sync *Unidirectional) iterate() (finished bool, err error) {
	if len(sync.directoryStack) == 0 && sync.currentDirectory == nil {
		finished = true
		return
	}

	if sync.currentDirectory == nil {
		sync.currentDirectory = sync.directoryStack[0]
		sync.directoryStack = sync.directoryStack[1:]
	}

	sync.currentlySyncing, err = sync.currentDirectory.NextEntry()

	if err != nil {
		return
	}

	if sync.currentlySyncing == nil {
		err = sync.currentDirectory.Close()
		sync.currentDirectory = nil
		return
	}

	switch entry := sync.currentlySyncing.(type) {

	case storage.DirectoryEntry:
		err = sync.syncDir(entry)

		if err != nil {
			return
		}

		var dir storage.Directory
		dir, err = entry.Open()

		if err != nil {
			return
		}

		sync.directoryStack = append(sync.directoryStack, dir)

	case storage.FileEntry:
		err = sync.syncFile(sync.currentDirectory.Entry(), entry)

		if err != nil {
			return
		}
	}

	return
}

func (sync *Unidirectional) run() {
	for true {
		finished, err := sync.iterate()

		// errors that happen during cleanup do not bubble as return values but are
		// instead logged on the object, so we deal with all errors this way
		sync.setLastErrorIfApplicable(err)

		if finished {
			sync.state = STATE_FINISHED
			sync.events <- &EventSyncFinished{}
			close(sync.events)
			return
		}

		if sync.lastError != nil {
			fmt.Printf("ERR %v\n", sync.lastError)
			sync.state = STATE_ERROR
			sync.events <- &EventError{err: sync.lastError}
			return
		}
	}
}

func (sync *Unidirectional) Start() (events chan EventInterface, err error) {
	if sync.state != STATE_NONE {
		err = errors.New("sync already started")
		return
	}

	sync.to, err = sync.config.To.Root()

	if err != nil {
		return
	}

	sync.from, err = sync.config.From.Root()

	if err != nil {
		_ = sync.from.Close()
		return
	}

	sync.directoryStack = append(sync.directoryStack, sync.from)

	events = make(chan EventInterface, 10)
	sync.events = events

	go sync.run()

	return
}

func (sync *Unidirectional) Resume() error {
	if sync.state != STATE_ERROR {
		return errors.New("sync not paused")
	}

	sync.lastError = nil

	go sync.run()

	return nil
}

func NewUnidirectionalSync(config Config) *Unidirectional {
	return &Unidirectional{
		config:         config,
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
		state:          STATE_NONE,
		directoryStack: make([]storage.Directory, 0, 10),
	}
}
