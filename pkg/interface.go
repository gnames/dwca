package dwca

import (
	"context"

	"github.com/gnames/dwca/pkg/config"
	"github.com/gnames/dwca/pkg/ent/eml"
	"github.com/gnames/dwca/pkg/ent/meta"
)

// Archive is an interface for Darwin Core Archive objects.
type Archive interface {
	// Config returns the configuration object of the archive.
	Config() config.Config

	// Meta returns the Meta object of the archive.
	Meta() *meta.Meta

	// EML returns the EML object of the archive.
	EML() *eml.EML

	// Load extracts the archive and loads data for EML and Meta.
	// Path determines internal location of the extracted archive.
	Load(path string) error

	// Close cleans up temporary files.
	Close() error

	// CoreSlice takes an offset and a limit and returns a slice of slices of
	// strings, each slice representing a row of the core file. If limit and
	// offset are provided, it returns the corresponding subset of the data.
	CoreSlice(offset, limit int) ([][]string, error)

	// CoreStream takes a channel and populates the channel with slices of
	// strings, each slice representing a row of the core file. The channel
	// is closed when the data is exhausted.
	CoreStream(context.Context, chan<- []string) error

	// ExtensionSlice takes an index, offset and limit and returns a slice of
	// slices of strings, each slice representing a row of the extension file.
	// Index corresponds the index of the extension in the extension list.
	// If limit and offset are provided, it returns the corresponding subset
	// of the data.
	ExtensionSlice(index, offset, limit int) ([][]string, error)

	// ExtensionStream takes an index and a channel and populates the channel
	// with slices of strings, each slice representing a row of the extension
	// file. The channel is closed when the data is exhausted.
	// Index corresponds the index of the extension in the extension list.
	ExtensionStream(ctx context.Context, index int, ch chan<- []string) error

	// Normalize creates a normalized version of Darwin Core Archive
	// with all known ambiguities resolved. The output is written to a file
	// with the provided fileName.
	Normalize() error

	// ZipNorgalized compresses a normalized version of Darwin Core Archive
	// to a ZIP file with the provided filePath.
	ZipNormalized(filePath string) error

	// TarGzNormalized compresses a normalized version of Darwin Core Archive
	// to a TAR file with the provided filePath.
	TarGzNormalized(filePath string) error
}
