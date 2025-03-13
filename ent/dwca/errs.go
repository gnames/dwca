package dwca

// ErrDownload should be used if download fails.
type ErrDownload struct {
	URL string
	Err error
}

func (e *ErrDownload) Error() string {
	return e.Err.Error()
}

// ErrExtractArchive should be used if DwCA file cannot be uncompressed.
type ErrExtractArchive struct {
	File string
	Err  error
}

func (e *ErrExtractArchive) Error() string {
	return e.Err.Error()
}

// ErrNoMetaFile is returned when the meta.xml file is not found in the
// extract directory.
type ErrNoMetaFile struct{}

func (e *ErrNoMetaFile) Error() string {
	return "meta.xml not found"
}

// ErrMultipleMetaFiles is returned when there are multiple meta.xml files in
// the extract directory.
type ErrMultipleMetaFiles struct{}

func (e *ErrMultipleMetaFiles) Error() string {
	return "multiple meta.xml files found"
}
