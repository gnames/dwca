package eml

import "fmt"

type ErrReader struct {
	OrigErr error
}

func (e *ErrReader) Error() string {
	return fmt.Sprintf("cannot read: %v", e.OrigErr)
}

type ErrDecoder struct {
	OrigErr error
}

func (e *ErrDecoder) Error() string {
	return fmt.Sprintf("cannot decode to EML: %v", e.OrigErr)
}
