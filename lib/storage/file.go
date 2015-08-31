package storage

import (
	"io"
)

type File interface {
	io.ReadCloser

	Entry() FileEntry
}
