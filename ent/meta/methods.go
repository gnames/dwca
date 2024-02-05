package meta

import (
	"bytes"
	"encoding/xml"
	"io"
)

// New reads an EML file from an io.Reader and returns an EML struct.
func New(r io.Reader) (*Meta, error) {
	bs, err := io.ReadAll(r)
	if err != nil {
		err = &ErrMetaReader{OrigErr: err}
		return nil, err
	}

	var res Meta
	decoder := xml.NewDecoder(bytes.NewReader(bs))
	err = decoder.Decode(&res)
	if err != nil {
		err = &ErrMetaDecoder{OrigErr: err}
		return nil, err
	}
	return &res, nil
}
