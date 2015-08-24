package sync

import (
	"github.com/DirtyHairy/csync/lib/storage"
)

type Config struct {
	From storage.Directory
	To   storage.Directory

	Verbose bool
	Dryrun  bool
}
