package dcfile

import (
	"fmt"
)

// ErrMetaFileNotFound is returned when the meta.xml file is not found in the
// extract directory.
type ErrMetaFileNotFound struct{}

func (e *ErrMetaFileNotFound) Error() string {
	return "meta.xml not found"
}

// ErrMultipleMetaFiles is returned when there are multiple meta.xml files in
// the extract directory.
type ErrMultipleMetaFiles struct{}

func (e *ErrMultipleMetaFiles) Error() string {
	return "multiple meta.xml files found"
}

// ErrUnknownArchiveType is returned when the file type is not supported.
type ErrUnknownArchiveType struct {
	FileType
}

func (e *ErrUnknownArchiveType) Error() string {
	return fmt.Sprintf("unknown file type: %s", e.FileType)
}

// ErrDir is returned when the directory is in an unknown state, or is not a
// directory.
type ErrDir struct {
	DirPath string
}

func (e *ErrDir) Error() string {
	return fmt.Sprintf("Directory '%s' is broken", e.DirPath)
}

// ErrFileNotFound is returned when the file is not found.
type ErrFileNotFound struct {
	Path string
}

func (e *ErrFileNotFound) Error() string {
	return fmt.Sprintf("file '%s' not found", e.Path)
}

// ErrExtract is returned when the extraction of the DwCA file fails.
type ErrExtract struct {
	Path string
	Err  error
}

func (e *ErrExtract) Error() string {
	return fmt.Sprintf("extracting '%s' failed: %v", e.Path, e.Err)
}

// ErrCoreRead is returned when reading the core file fails.
type ErrCoreRead struct {
	Err error
}

func (e *ErrCoreRead) Error() string {
	return fmt.Sprintf("reading core file failed: %v", e.Err)
}

type ErrExtensionRead struct {
	Err error
}

// ErrExtensionRead is returned when reading the extension file fails.
func (e *ErrExtensionRead) Error() string {
	return fmt.Sprintf("reading extension file failed: %v", e.Err)
}
