package local

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/DirtyHairy/csync/lib/storage"
)

func getFSInstace() storage.StorageProvider {
	fs, err := NewLocalFS("./test_artifacts")

	if err != nil {
		panic(err)
	}

	return fs
}

func getFSRoot() storage.Directory {
	root, err := getFSInstace().Root()

	if err != nil {
		panic(err)
	}

	return root
}

func getTempFSInstance() (storage.StorageProvider, error) {
	path, err := ioutil.TempDir("", "csync_test")

	if err != nil {
		return nil, err
	}

	fs, err := NewLocalFS(path)

	if err != nil {
		return nil, err
	}

	return fs, nil
}

func getTempFSRoot() (storage.Directory, error) {
	fs, err := getTempFSInstance()

	if err != nil {
		return nil, err
	}

	root, err := fs.Root()

	if err != nil {
		return nil, err
	}

	return root, nil
}

func destroyTempFS(fs storage.Directory) {
	directory, ok := fs.(*directory)

	if !ok {
		panic("no an instance of directory")
	}

	if err := os.RemoveAll(directory.Entry().(*directoryEntry).realPath()); err != nil {
		panic(err)
	}
}

func checkDirectoryContents(directory storage.Directory, expectedContents []string) (entries map[string]storage.Entry, e error) {
	entries = make(map[string]storage.Entry)
	e = nil

	contents := make([]string, 0, 10)

	for entry, err := directory.NextEntry(); entry != nil; entry, err = directory.NextEntry() {
		if err != nil {
			e = errors.New(fmt.Sprintf("error while iterating over directory entries: %v", err))
			return
		}

		contents = append(contents, entry.Name())
		entries[entry.Name()] = entry
	}

	if len(contents) != len(expectedContents) {
		e = errors.New(fmt.Sprintf("expected %d dir entries, got %d instead", len(expectedContents), len(contents)))
		return
	}

	sort.Sort(sort.StringSlice(contents))
	sort.Sort(sort.StringSlice(expectedContents))

	for idx, filename := range contents {
		if filename != expectedContents[idx] {
			e = errors.New(fmt.Sprintf("directory listing differs at %d: expected %s, got %s", idx, expectedContents[idx], filename))
			return
		}
	}

	return
}

func expectedRootEntries() []string {
	return []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "foo", "bar"}
}
