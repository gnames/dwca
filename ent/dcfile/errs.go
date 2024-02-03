package dcfile

import "fmt"

// ErrUnknownFileType is returned when the file type is not supported.
type ErrUnknownFileType struct {
	FileType
}

func (e ErrUnknownFileType) Error() string {
	return fmt.Sprintf("unknown file type: %s", e.FileType)
}

// ErrDir is returned when the directory is in an unknown state, or is not a
// directory.
type ErrDir struct {
	DirPath string
}

func (e ErrDir) Error() string {
	return fmt.Sprintf("Directory '%s' is broken", e.DirPath)
}

// ErrFileNotFound is returned when the file is not found.
type ErrFileNotFound struct {
	Path string
}

func (e ErrFileNotFound) Error() string {
	return fmt.Sprintf("file '%s' not found", e.Path)
}

type ErrExtract struct {
	Path string
	Err  error
}

func (e ErrExtract) Error() string {
	return fmt.Sprintf("extracting '%s' failed: %v", e.Path, e.Err)
}
