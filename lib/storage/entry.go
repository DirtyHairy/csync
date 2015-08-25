package storage

import (
	"time"
)

type Entry interface {
	Name() string
	Path() string
	Mtime() time.Time
	SetMtime(time.Time) error

	Remove() error
	Rename(newname string) (Entry, error)
}
