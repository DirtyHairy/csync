package config

type Config interface {
	Type() string
	Path() string
}

type MutableConfig interface {
	Config

	SetPath(string)
}
