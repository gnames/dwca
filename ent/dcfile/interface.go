package dcfile

// DCFile represents a Darwin Core Archive file. It is normally a compressed
// tar file or zip file with a set of files inside that correspond to DwCA
// structure and format.
type DCFile interface {
	//  Init cleans up or delets temporary directories.
	Init() error

	//  Extract extracts the content of the DwCA file to a temporary directory.
	Extract() error

	// Close removes the temporary directory with the extracted content.
	Close() error
}
