package dcfile

import (
	"context"

	"github.com/gnames/dwca/pkg/ent/meta"
)

// DCFile represents a Darwin Core Archive file. It is normally a compressed
// tar file or zip file with a set of files inside that correspond to DwCA
// structure and format.
type DCFile interface {
	//  ResetTempDirs cleans up or deletes temporary directories.
	ResetTempDirs() error

	//	SetFilePath sets the path to the DwCA file.
	SetFilePath(string)

	//  Extract extracts the content of the DwCA file to a temporary directory.
	Extract() error

	// ArchiveDir returns the path to the temporary directory
	// where DwCA data is located.
	ArchiveDir(path string) (string, error)

	// CoreData returns the content of the core file as a slice of slices of
	// strings. Each slice of strings represents a row in the core file.
	CoreData(root string, meta *meta.Meta, offset, limit int) ([][]string, error)

	// CoreStream populates a channel that streams the content of the core file
	// as a slice of strings. Each slice of strings represents a row in the core
	// file.
	CoreStream(
		ctx context.Context,
		root string,
		meta *meta.Meta,
		coreChan chan<- []string,
	) (int, error)

	// ExtensionData returns the content of the extension file as a slice of
	// slices of strings. Each slice of strings represents a row in the extension
	// file. The index is the index of the extension file in the meta file.
	// The offset and limit are used to paginate the results.
	ExtensionData(
		index int, root string,
		meta *meta.Meta,
		offset, limit int,
	) ([][]string, error)

	// ExtensionStream populates a channel that streams the content of the
	// extension file as a slice of strings. Each slice of strings represents a
	// row in the extension file. The index is the index of the extension file in
	// the meta file.
	ExtensionStream(
		ctx context.Context,
		index int, root string,
		meta *meta.Meta,
		extChan chan<- []string,
	) (int, error)

	// ExportCSVStream saves the content of a stream to a file. The file is a
	// comma-separated file with the first row being the header. The header is
	// defined by the fields parameter. This function is used to export Core or
	// Extension data to a file.
	ExportCSVStream(
		ctx context.Context,
		file string,
		headers []string,
		delim string,
		outChan <-chan []string) error

	// SaveToFile saves bytes slice to a file with the provided name.
	SaveToFile(fileName string, bs []byte) error

	// Zip compresses the content of the temporary output directory to a
	// ZIP file with the provided filePath.
	Zip(inputDir, zipFile string) error

	// TarGz compresses the content of the temporary output directory to a
	// TAR file with the provided filePath.
	TarGz(inputDir, tarFile string) error

	// Close removes the temporary directory with the extracted content.
	Close() error
}
