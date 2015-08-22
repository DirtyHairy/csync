package storage

import (
	"io"
)

type WritableFile interface {
	File
	io.Writer
}
