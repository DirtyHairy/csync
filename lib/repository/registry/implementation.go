package registry

import (
	"github.com/DirtyHairy/csync/lib/repository"
	"github.com/DirtyHairy/csync/lib/repository/registry/config"
)

type registry map[string](repository.Repository)

func (r registry) All() map[string]repository.Repository {
	clone := make(map[string]repository.Repository)

	for key, value := range r {
		clone[key] = value
	}

	return clone
}

func (r registry) Keys() []string {
	keys := make([]string, len(r))

	i := 0
	for key := range r {
		keys[i] = key
		i++
	}

	return keys
}

func (r registry) Size() int {
	return len(r)
}

func (r registry) Get(key string) repository.Repository {
	return r[key]
}

func (r registry) Set(key string, repo repository.Repository) {
	r[key] = repo
}

func (r registry) Unset(key string) {
	delete(r, key)
}

func (r registry) Marshal() config.Config {
	cfg := config.NewMutableConfig()

	for key, repo := range r {
		cfg.AddConfig(key, repo.Marshal())
	}

	return cfg
}

func NewRegistry() MutableRegistry {
	return registry(make(map[string]repository.Repository))
}
