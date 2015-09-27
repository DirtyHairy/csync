package environment

import (
	"os"

	"github.com/DirtyHairy/csync/lib/environment/config"
	"github.com/blang/semver"
)

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
	} else if err != nil {
		return err
	}

	_ = config

	return nil
}

func (e *environment) Save() error {
	if !e.dirty {
		return nil
	}

	locator := config.NewLocator()
	manager := config.NewManager(locator)

	config := config.NewConfig(e.Version())
	err := manager.Save(config)

	if err != nil {
		return err
	}

	e.dirty = false
	return nil
}

func New() MutableEnvironment {
	return &environment{}
}
