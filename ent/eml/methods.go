package eml

import (
	"bytes"
	"encoding/xml"
	"io"
)

// New reads an EML file from an io.Reader and returns an EML struct.
func New(r io.Reader) (*EML, error) {
	bs, err := io.ReadAll(r)
	if err != nil {
		err = &ErrReader{OrigErr: err}
		return nil, err
	}

	var res EML
	decoder := xml.NewDecoder(bytes.NewReader(bs))
	err = decoder.Decode(&res)
	if err != nil {
		err = &ErrDecoder{OrigErr: err}
		return nil, err
	}
	return &res, nil
}
