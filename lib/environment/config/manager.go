package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/blang/semver"
)

const CONFIG_FILE_NAME = "config.json"

type Manager interface {
	Load() (Config, error)
	Save(Config) error
}

type manager struct {
	locator ConfigLocator
}

func unmarshalConfig(jsonData []byte) (config Config, err error) {
	var proto_config struct {
		CsyncVersion *semver.Version `json:"csyncVersion"`
	}

	err = json.Unmarshal(jsonData, &proto_config)

	if err != nil {
		return
	}

	if proto_config.CsyncVersion == nil {
		err = errors.New("csyncVersion missing")
		return
	}

	config = newConfig(*proto_config.CsyncVersion)

	err = json.Unmarshal(jsonData, &config)

	return
}

func (m manager) Load() (config Config, err error) {
	directory, err := m.locator.Locate()

	if err != nil {
		return
	}

	fileContents, err := ioutil.ReadFile(filepath.Join(directory, CONFIG_FILE_NAME))

	if err != nil {
		return
	}

	return unmarshalConfig(fileContents)
}

func marshalConfig(config Config) (jsonData []byte, err error) {
	return json.MarshalIndent(config, "", "  ")
}

func (m manager) Save(config Config) (err error) {
	directory, err := m.locator.Locate()

	if err != nil {
		return
	}

	jsonData, err := marshalConfig(config)

	if err != nil {
		return
	}

	return ioutil.WriteFile(filepath.Join(directory, CONFIG_FILE_NAME), jsonData, 0660)
}

func NewManager(locator ConfigLocator) Manager {
	return manager{
		locator: locator,
	}
}
