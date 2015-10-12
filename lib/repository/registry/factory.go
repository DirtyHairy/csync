package registry

import (
	"github.com/DirtyHairy/csync/lib/repository"
	"github.com/DirtyHairy/csync/lib/repository/registry/config"
)

func CreateFromConfig(cfg config.Config) (registry MutableRegistry, err error) {
	registry = NewRegistry()

	for key, value := range cfg.RepositoryConfigs() {
		var repo repository.Repository

		repo, err = repository.CreateFromConfig(value)

		if err != nil {
			return
		}

		registry.Set(key, repo)
	}

	return
}
