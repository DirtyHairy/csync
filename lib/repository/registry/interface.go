package registry

import (
	"github.com/DirtyHairy/csync/lib/repository"
	"github.com/DirtyHairy/csync/lib/repository/registry/config"
)

type Registry interface {
	All() map[string]repository.Repository
	Keys() []string
	Get(key string) repository.Repository

	Marshal() config.Config
}

type MutableRegistry interface {
	Registry

	Set(key string, repo repository.Repository)
	Unset(key string)
}
