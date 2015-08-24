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

func getFSRoot(path string) storage.Directory {
	fs, err := local.NewLocalFS(path)

	if err != nil {
		fmt.Printf("unable to open %s\n", path)

		os.Exit(1)
	}

	root, err := fs.Root()

	if err != nil {
		panic(err)
	}

	return root
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	from := getFSRoot(os.Args[1])
	to := getFSRoot(os.Args[2])

	syncConfig := sync.Config{
		From: from,
		To:   to,
	}

	usync := sync.NewUnidirectionalSync(syncConfig)

	if err := usync.Execute(); err != nil {
		fmt.Printf("\nSYNC FAILED: %v\n", err)
	} else {
		fmt.Println("\nSYNC SUCCESS")
	}
}
