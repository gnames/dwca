package meta

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
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

func (m *Meta) ToData() *Data {
	data := &Data{}
	data.ExtensionsData = make(map[string]ExtensionData)
	data.CoreData = m.Core.toCoreData()
	for _, ext := range m.Extensions {
		file := filepath.Base(ext.Files.Location)
		name := stripExt(file)
		data.ExtensionsData[name] = ext.toExtensionData()
	}
	return data
}

func (c *Core) toCoreData() CoreData {
	var term string
	if c.ID.Term != "" {
		term = c.ID.Term
	}
	var idx int
	idxRes, err := strconv.Atoi(c.ID.Index)
	if err == nil {
		idx = idxRes
	}
	coreData := CoreData{
		Index:      idx,
		TermFull:   term,
		Term:       filepath.Base(term),
		FieldsData: make(map[string]FieldData),
	}
	for _, field := range c.Fields {
		idx = 0
		idxRes, err := strconv.Atoi(field.Index)

		if err == nil {
			idx = idxRes
			fmt.Printf("TRM: %#v\n", coreData.Term)
			if idx == 0 && coreData.Term == "" {
				fmt.Printf("FLD: %#v\n", field)
				coreData.Term = field.Term
				coreData.Term = filepath.Base(field.Term)
			}
		}

		term := filepath.Base(field.Term)
		coreData.FieldsData[term] = FieldData{
			Index:    idx,
			TermFull: field.Term,
			Term:     term,
		}
	}
	return coreData
}

func (e *Extension) toExtensionData() ExtensionData {
	idx := 0
	idxRes, err := strconv.Atoi(e.CoreID.Index)
	if err == nil {
		idx = idxRes
	}
	extData := ExtensionData{
		CoreIndex:  idx,
		FieldsData: make(map[string]FieldData),
	}
	for _, field := range e.Fields {
		term := filepath.Base(field.Term)
		idx = 0
		idxRes, err := strconv.Atoi(field.Index)
		if err == nil {
			idx = idxRes
		}
		extData.FieldsData[term] = FieldData{
			Index:    idx,
			TermFull: field.Term,
			Term:     term,
		}
	}
	return extData
}

func stripExt(filename string) string {
	ext := len(filepath.Ext(filename))
	end := len(filename) - ext
	return filename[:end]
}
