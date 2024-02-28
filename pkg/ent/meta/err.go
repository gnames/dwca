package meta

import "fmt"

// ErrMetaReader is an error type for reading meta.xml files.
type ErrMetaReader struct {
	// OrigErr is the original error.
	OrigErr error
}

func (e *ErrMetaReader) Error() string {
	return fmt.Sprintf("cannot read: %v", e.OrigErr)
}

// ErrMetaDecoder is an error type for decoding meta.xml files.
type ErrMetaDecoder struct {
	//  OrigErr is the original error.
	OrigErr error
}

func (e *ErrMetaDecoder) Error() string {
	return fmt.Sprintf("cannot decode meta.xml: %v", e.OrigErr)
}
