package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Path         string
	DownloadPath string
	ExtractPath  string
	WithCleanup  bool
}

type Option func(*Config)

func OptPath(s string) Option {
	return func(c *Config) {
		c.Path = s
	}
}

func OptWithCleanup(b bool) Option {
	return func(c *Config) {
		c.WithCleanup = b
	}
}

func New(opts ...Option) Config {
	path, err := os.UserCacheDir()
	if err != nil {
		path = os.TempDir()
	}

	path = filepath.Join(path, "dwca_go")
	c := Config{Path: path}
	for _, opt := range opts {
		opt(&c)
	}
	c.DownloadPath = filepath.Join(c.Path, "download")
	c.ExtractPath = filepath.Join(c.Path, "extract")
	return c
}
