package config

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

type ConfigLocator interface {
	Locate() (string, error)
}

type defaultLocator struct{}

func (defaultLocator) Locate() (directory string, err error) {
	var parentDirectory string

	if runtime.GOOS == "windows" {
		parentDirectory = os.Getenv("APPDATA")
	}

	if parentDirectory == "" {
		var currentUser *user.User

		currentUser, err = user.Current()

		if err != nil {
			return
		}

		parentDirectory = currentUser.HomeDir
	}

	directory = filepath.Join(parentDirectory, ".csync")

	_, err = os.Stat(directory)

	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(directory, 0755)
	}

	return
}

func NewLocator() ConfigLocator {
	return defaultLocator{}
}
