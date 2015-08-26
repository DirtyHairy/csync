package sync

import (
	"fmt"

	"github.com/DirtyHairy/csync/lib/storage"
)

// Start syncing a file or creating a directory
type EventStartSyncing struct {
	entry storage.Entry
}

func (e *EventStartSyncing) Description() string {
	return fmt.Sprintf("syncing %s", e.entry.Path())
}

func (e *EventStartSyncing) Entry() storage.Entry {
	return e.entry
}

// The sync process encountered an error
type EventError struct {
	err error
}

func (e *EventError) Description() string {
	return fmt.Sprintf("sync failed: %v", e.err)
}

func (e *EventError) Error() error {
	return e.err
}

// Syncing finished
type EventSyncFinished struct{}

func (e *EventSyncFinished) Description() string {
	return "sync complete"
}
