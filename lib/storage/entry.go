package storage

import (
	"time"
)

type Entry interface {
	Name() string
	Path() string
	Mtime() time.Time
}
