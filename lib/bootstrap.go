package lib

import (
	"github.com/DirtyHairy/csync/lib/environment"
	"github.com/blang/semver"
)

func Bootstrap(e interface {
	environment.VersionReceiver
}) {
	e.SetVersion(semver.MustParse("0.0.1"))
}
