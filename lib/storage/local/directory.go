package local

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/DirtyHairy/csync/lib/storage"
)

type directory struct {
	entry *directoryEntry

	file *os.File

	directorySlice []storage.Entry
	idx            int

	closed bool
}

func (d *directory) closedError() error {
	return errors.New("directory already closed")
}

func (d *directory) Entry() storage.DirectoryEntry {
	return d.entry
}

func (d *directory) Stat(path string) (storage.Entry, error) {
	if d.closed {
		return nil, d.closedError()
	}

	realPath := filepath.Join(d.entry.realPath(), filepath.FromSlash(path))
	fileInfo, err := statLinkTarget(realPath)

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	newPath := filepath.Join(d.entry.path, filepath.FromSlash(path))

	if fileInfo.IsDir() {
		return newDirectoryEntry(newPath, d.entry.prefix, fileInfo)
	} else {
		return newFileEntry(newPath, d.entry.prefix, fileInfo)
	}
}

func (d *directory) readSlice() error {
	var err error

	if d.file == nil {
		d.file, err = os.Open(d.entry.realPath())

		if err != nil {
			return err
		}
	}

	fi, err := d.file.Readdir(10)

	if err != nil && err != io.EOF {
		return err
	}

	slice := d.directorySlice

	if slice == nil || len(slice) != len(fi) {
		slice = make([]storage.Entry, len(fi))
	}

	for i, fi := range fi {
		var (
			entry storage.Entry
			err   error
		)

		path := filepath.Join(d.entry.path, fi.Name())

		if fi.Mode()&os.ModeSymlink != 0 {
			fi, err = statLinkTarget(filepath.Join(d.entry.realPath(), fi.Name()))

			if err != nil {
				return err
			}
		}

		if fi.IsDir() {
			entry, err = newDirectoryEntry(path, d.entry.prefix, fi)
		} else {
			entry, err = newFileEntry(path, d.entry.prefix, fi)
		}

		if err != nil {
			return err
		}

		slice[i] = entry
	}

	d.idx = 0
	d.directorySlice = slice

	return nil
}

func (d *directory) NextEntry() (storage.Entry, error) {
	if d.closed {
		return nil, d.closedError()
	}

	if d.directorySlice == nil || (d.idx >= len(d.directorySlice) && len(d.directorySlice) != 0) {
		err := d.readSlice()

		if err != nil {
			return nil, err
		}
	}

	if len(d.directorySlice) == 0 {
		return nil, nil
	}

	entry := d.directorySlice[d.idx]
	d.idx++

	return entry, nil
}

func (d *directory) Rewind() error {
	var err error

	if d.closed {
		return d.closedError()
	}

	if d.file != nil {
		err = d.file.Close()

		if err != nil {
			return err
		}
	}

	d.file = nil
	d.idx = 0
	d.directorySlice = nil

	return nil
}

func (d *directory) CreateFile(path string) (storage.WritableFile, error) {
	path = filepath.FromSlash(path)

	realpath := filepath.Join(d.entry.realPath(), path)

	file, err := os.OpenFile(realpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0644)

	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(realpath)

	cleanup := func() {
		_ = file.Close()
		_ = os.Remove(realpath)
	}

	if err != nil {
		cleanup()
		return nil, err
	}

	entry, err := newFileEntry(filepath.Join(d.entry.path, path), d.entry.prefix, fileInfo)

	if err != nil {
		cleanup()
		return nil, err
	}

	return newFile(entry, file), nil
}

func (d *directory) Mkdir(path string) (storage.Directory, error) {
	var err error

	path = filepath.FromSlash(path)
	realpath := filepath.Join(d.entry.prefix, path)

	err = os.MkdirAll(realpath, 0755)

	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(realpath)

	if err != nil {
		return nil, err
	}

	entry, err := newDirectoryEntry(filepath.Join(d.entry.path, path), d.entry.prefix, fileInfo)

	if err != nil {
		return nil, err
	}

	return newDirectory(entry), nil
}

func (d *directory) Close() error {
	if d.closed {
		return nil
	}

	err := d.Rewind()

	d.closed = true

	return err
}

func newDirectory(entry *directoryEntry) *directory {
	return &directory{
		entry: entry,
	}
}
