package push

import (
	"fmt"

	"github.com/DirtyHairy/csync/lib/storage"
	"github.com/DirtyHairy/csync/lib/storage/local"
	"github.com/DirtyHairy/csync/lib/sync"
)

func Execute(config Config) error {
	var err error

	from, err := local.NewLocalFS(config.SourceRepoId)

	if err != nil {
		return err
	}

	to, err := local.NewLocalFS(config.TargetRepoId)

	if err != nil {
		return err
	}

	syncInstance := sync.NewUnidirectionalSync(sync.Config{
		From: from,
		To:   to,
	})

	events, err := syncInstance.Start()

	if err != nil {
		return err
	}

EventLoop:
	for event := range events {
		switch event := event.(type) {

		case *sync.EventStartSyncing:
			switch entry := event.Entry().(type) {

			case storage.DirectoryEntry:
				fmt.Printf("creating %s\n", entry.Path())

			case storage.FileEntry:
				fmt.Printf("syncing %s\n", entry.Path())

			}

		case *sync.EventError:
			fmt.Printf("\nSYNC FAILED: %v\n", event.Error())
			break EventLoop

		case *sync.EventSyncFinished:
			fmt.Println("\nSYNC SUCCESSFUL\n")
			break EventLoop
		}
	}

	return nil
}
