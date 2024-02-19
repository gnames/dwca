package dcfile

import (
	"context"

	"github.com/gnames/dwca/ent/meta"
)

// DCFile represents a Darwin Core Archive file. It is normally a compressed
// tar file or zip file with a set of files inside that correspond to DwCA
// structure and format.
type DCFile interface {
	//  Init cleans up or deletes temporary directories.
	Init() error

	//  Extract extracts the content of the DwCA file to a temporary directory.
	Extract() error

	// ArchiveDir returns the path to the temporary directory
	// where DwCA data is located.
	ArchiveDir() (string, error)

	// CoreData returns the content of the core file as a slice of slices of
	// strings. Each slice of strings represents a row in the core file.
	CoreData(meta *meta.Meta, offset, limit int) ([][]string, error)

	// CoreStream populates a channel that streams the content of the core file
	// as a slice of strings. Each slice of strings represents a row in the core
	// file.
	CoreStream(
		ctx context.Context,
		meta *meta.Meta,
		coreChan chan<- []string,
	) error

	// ExtensionData returns the content of the extension file as a slice of
	// slices of strings. Each slice of strings represents a row in the extension
	// file. The index is the index of the extension file in the meta file.
	// The offset and limit are used to paginate the results.
	ExtensionData(
		index int, meta *meta.Meta,
		offset, limit int,
	) ([][]string, error)

	// ExtensionStream populates a channel that streams the content of the
	// extension file as a slice of strings. Each slice of strings represents a
	// row in the extension file. The index is the index of the extension file in
	// the meta file.
	ExtensionStream(
		ctx context.Context,
		index int,
		meta *meta.Meta,
		extChan chan<- []string,
	) error

	// ExportCSVStream saves the content of a stream to a file. The file is a
	// comma-separated file with the first row being the header. The header is
	// defined by the fields parameter. This function is used to export Core or
	// Extension data to a file.
	ExportCSVStream(
		ctx context.Context,
		file string,
		fields []string,
		outChan <-chan []string) error

	// SaveToFile saves bytes slice to a file with the provided name.
	SaveToFile(fileName string, bs []byte) error

	// ZipOutput compresses the content of the temporary output directory to a
	// ZIP file with the provided filePath.
	ZipOutput(filePath string) error

	// TarGzOutput compresses the content of the temporary output directory to a
	// TAR file with the provided filePath.
	TarGzOutput(filePath string) error

	// Close removes the temporary directory with the extracted content.
	Close() error
}
