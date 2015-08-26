package main

import (
	"fmt"
	"os"

	"github.com/DirtyHairy/csync/lib/storage"
	"github.com/DirtyHairy/csync/lib/storage/local"
	"github.com/DirtyHairy/csync/lib/sync"
)

func usage() {
	fmt.Println("usage: csync repo1 repo2")
	os.Exit(1)
}

func getFS(path string) storage.StorageProvider {
	fs, err := local.NewLocalFS(path)

	if err != nil {
		fmt.Printf("unable to open %s\n", path)

		os.Exit(1)
	}

	return fs
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	from := getFS(os.Args[1])
	to := getFS(os.Args[2])

	syncConfig := sync.Config{
		From: from,
		To:   to,
	}

	syncInstance := sync.NewUnidirectionalSync(syncConfig)

	events, err := syncInstance.Start()

	if err != nil {
		panic(err)
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
}
