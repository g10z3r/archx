package anthill

import "time"

var ignoredMap = map[string]struct{}{
	".git":    {},
	".docker": {},

	".vscode":  {},
	".idea":    {},
	".eclipse": {},

	"dist":    {},
	"docker":  {},
	"assets":  {},
	"vendor":  {},
	"build":   {},
	"scripts": {},
	"ci":      {},
	"log":     {},
	"logs":    {},
}

type ConfigOption func(*Config)

type Config struct {
	ignoredList map[string]struct{}
	timeout     time.Duration
	selectedDir string
}

func WithIgnoredList(list ...string) ConfigOption {
	return func(c *Config) {
		for _, item := range list {
			c.ignoredList[item] = struct{}{}
		}
	}
}

func WithTimeout(timeout time.Duration) ConfigOption {
	return func(c *Config) {
		c.timeout = timeout
	}
}

func WithSelectedDir(dir string) ConfigOption {
	return func(c *Config) {
		c.selectedDir = dir
	}
}

func NewConfig(options ...ConfigOption) *Config {
	config := &Config{}
	for _, option := range options {
		option(config)
	}

	return config
}

func DefaultConfig(options ...ConfigOption) *Config {
	config := &Config{
		ignoredList: ignoredMap,
	}

	for _, option := range options {
		option(config)
	}

	return config
}
