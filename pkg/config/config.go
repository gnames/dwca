package config

import (
	"os"
	"path/filepath"
)

// Config is a configuration object for the Darwin Core Archive (DwCA)
// data processing.
type Config struct {
	// RootPath is the root path for all temporary files.
	RootPath string

	// DownloadPath is used to store downloaded files.
	DownloadPath string

	// ExtractPath is used to store extracted files of DwCA archive.
	ExtractPath string

	// OutputPath is used to store uncompressed files of a normalized
	// DwCA archive. This files are created from the original DwCA archive
	// data.
	OutputPath string

	// JobsNum is the number of concurrent jobs to run.
	JobsNum int
}

// Option is a function type that allows to standardize how options to
// the configuration are organized.
type Option func(*Config)

// OptPath sets the root path for all temporary files.
func OptPath(s string) Option {
	return func(c *Config) {
		c.RootPath = s
	}
}

// New creates a new Config object with default values, and allows to
// override them with options.
func New(opts ...Option) Config {
	path, err := os.UserCacheDir()
	if err != nil {
		path = os.TempDir()
	}

	path = filepath.Join(path, "dwca_go")
	c := Config{
		RootPath: path,
		JobsNum:  5,
	}

	for _, opt := range opts {
		opt(&c)
	}

	c.DownloadPath = filepath.Join(c.RootPath, "download")
	c.ExtractPath = filepath.Join(c.RootPath, "extract")
	c.OutputPath = filepath.Join(c.RootPath, "output")
	return c
}
