package environment

import (
	"errors"
	"fmt"
	"os"

	"github.com/DirtyHairy/csync/lib/environment/config"
	"github.com/blang/semver"
)

const CSYNC_VERSION = "0.0.1"

type environment struct {
	version semver.Version

	dirty bool
}

func (e *environment) Version() semver.Version {
	return e.version
}

func (e *environment) SetVersion(version semver.Version) {
	e.version = version
}

func (e *environment) Load() error {
	locator := config.NewLocator()
	manager := config.NewManager(locator)

	config, err := manager.Load()

	if err != nil && os.IsNotExist(err) {
		e.dirty = true
		return nil
	}

	if err != nil {
		return err
	}

	loadedVersion := config.CsyncVersion()

	if loadedVersion.LT(e.version) {
		e.dirty = true
	}

	if loadedVersion.GT(e.version) {
		return errors.New(fmt.Sprintf(
			"config was written by newer csync version (%s vs. %s)", loadedVersion, e.version))
	}

	return nil
}

func (e *environment) Save() error {
	if !e.dirty {
		return nil
	}

	locator := config.NewLocator()
	manager := config.NewManager(locator)

	config := config.NewMutableConfig(e.Version())
	err := manager.Save(config)

	if err != nil {
		return err
	}

	e.dirty = false
	return nil
}

func New() MutableEnvironment {
	return &environment{
		version: semver.MustParse(CSYNC_VERSION),
	}
}
