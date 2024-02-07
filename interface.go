package dwca

import (
	"github.com/gnames/dwca/config"
	"github.com/gnames/dwca/ent/eml"
	"github.com/gnames/dwca/ent/meta"
)

// Archive is an interface for Darwin Core Archive objects.
type Archive interface {
	// Config returns the configuration object of the archive.
	Config() config.Config

	// Load extracts the archive and loads data for EML and Meta.
	Load() error

	// Meta returns the Meta object of the archive.
	Meta() *meta.Meta

	// EML returns the EML object of the archive.
	EML() *eml.EML

	// CoreData takes an offset and a limit and returns a slice of slices of
	// strings, each slice representing a row of the core file. If limit and
	// offset are provided, it returns the corresponding subset of the data.
	CoreData(offset, limit int) ([][]string, error)

	// CoreStream takes a channel and populates the channel with slices of
	// strings, each slice representing a row of the core file. The channel
	// is closed when the data is exhausted.
	CoreStream(chan<- []string) error

	// ExtensionData takes an index, offset and limit and returns a slice of
	// slices of strings, each slice representing a row of the extension file.
	// Index corresponds the index of the extension in the extension list.
	// If limit and offset are provided, it returns the corresponding subset
	// of the data.
	ExtensionData(index, offset, limit int) ([][]string, error)

	// ExtensionStream takes an index and a channel and populates the channel
	// with slices of strings, each slice representing a row of the extension
	// file. The channel is closed when the data is exhausted.
	// Index corresponds the index of the extension in the extension list.
	ExtensionStream(index int, ch chan<- []string) error
}
