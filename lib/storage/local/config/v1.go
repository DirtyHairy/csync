package config

type config_v1 struct {
	Type_ string `json:"type"`
	Path_ string `json:"path"`
}

func (c *config_v1) Type() string {
	return c.Type_
}

func (c *config_v1) Path() string {
	return c.Path_
}

func (c *config_v1) SetPath(path string) {
	c.Path_ = path
}
