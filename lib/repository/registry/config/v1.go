package config

import (
	"encoding/json"

	"github.com/blang/semver"

	"github.com/DirtyHairy/csync/lib/repository/local_config"
)

type config_v1 struct {
	forVersion semver.Version
	configs    map[string]local_config.Config
}

func (c *config_v1) RepositoryConfigs() map[string]local_config.Config {
	clone := make(map[string]local_config.Config)

	if c.configs != nil {
		for key, value := range c.configs {
			clone[key] = value
		}
	}

	return clone
}

func (c *config_v1) AddConfig(key string, cfg local_config.Config) {
	if c.configs == nil {
		c.configs = make(map[string]local_config.Config)
	}

	c.configs[key] = cfg
}

func (c *config_v1) UnmarshalJSON(jsonData []byte) (err error) {
	c.configs = make(map[string]local_config.Config)

	stage1 := make(map[string]json.RawMessage)

	err = json.Unmarshal(jsonData, &stage1)
	if err != nil {
		return
	}

	for key, rawData := range stage1 {
		target := local_config.NewConfig(c.forVersion)
		err = json.Unmarshal(rawData, &target)

		if err != nil {
			return
		}

		c.configs[key] = target
	}

	return
}

func (c *config_v1) MarshalJSON() (jsonData []byte, err error) {
	return json.Marshal(c.configs)
}
