package storage

import (
	"io"
)

type File interface {
	io.Reader

	Entry() FileEntry

	Close() error
}
