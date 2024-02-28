package config

import (
	"log/slog"
	"os"
	"path/filepath"
)

var (
	outputCompression = "zip"
	outputCSVType     = "tsv"
	jobsNum           = 5
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

	// OutputArchiveCompression is the compression format to use when
	// creating the output archive. It can be "zip" or "tar.gz".
	OutputArchiveCompression string

	// OutputCSVType is the type of CSV files. Can be "csv" or "tsv"
	OutputCSVType string

	// JobsNum is the number of concurrent jobs to run.
	JobsNum int
}

// Option is a function type that allows to standardize how options to
// the configuration are organized.
type Option func(*Config)

// OptRootPath sets the root path for all temporary files.
func OptRootPath(s string) Option {
	return func(c *Config) {
		c.RootPath = s
	}
}

// OptOutputArchiveCompression sets the compression format to use when
// creating the output archive. It can be "zip" or "tar.gz".
func OptArchiveCompression(s string) Option {
	return func(c *Config) {
		if s != "zip" && s != "tar" {
			slog.Warn(
				"Entered compression format is not supported. Using default format",
				"input", s, "default", outputCompression,
			)
			s = outputCompression
		}
		c.OutputArchiveCompression = s
	}
}

// OptOutputCSVType sets the type of CSV files. Can be "csv" or "tsv"
func OptOutputCSVType(s string) Option {
	return func(c *Config) {
		if s != "csv" && s != "tsv" {
			slog.Warn(
				"Entered CSV type is not supported. Using default format",
				"bad-input", s, "default", outputCSVType,
			)
			s = outputCSVType
		}
		c.OutputCSVType = s
	}
}

// OptJobsNum sets the number of concurrent jobs to run.
func OptJobsNum(n int) Option {
	return func(c *Config) {
		if n < 1 || n > 100 {
			slog.Warn(
				"Unsupported number of jobs (supported: 1-100). Using default value",
				"bad-input", n, "default", jobsNum,
			)
			n = jobsNum
		}
		c.JobsNum = n
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
		RootPath:                 path,
		OutputArchiveCompression: outputCompression,
		OutputCSVType:            outputCSVType,
		JobsNum:                  jobsNum,
	}

	for _, opt := range opts {
		opt(&c)
	}

	c.DownloadPath = filepath.Join(c.RootPath, "download")
	c.ExtractPath = filepath.Join(c.RootPath, "extract")
	c.OutputPath = filepath.Join(c.RootPath, "output")
	return c
}
