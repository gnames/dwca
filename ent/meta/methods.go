package meta

import (
	"bytes"
	"encoding/xml"
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

// New reads an Meta file from an io.Reader and returns an EML struct.
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

func (m *Meta) Simplify() *MetaSimple {
	data := &MetaSimple{}
	data.ExtensionsData = make(map[string]ExtensionData)
	data.CoreData = m.Core.toCoreData()

	for _, ext := range m.Extensions {
		file := filepath.Base(ext.Files.Location)
		name := stripExt(file)
		if ext.RowType != "" {
			name = filepath.Base(ext.RowType)
		}
		name = strings.ToLower(name)
		data.ExtensionsData[name] = ext.toExtensionData()
	}
	return data
}

func (c *Core) toCoreData() CoreData {
	var termFull string
	if c.ID.Term != "" {
		termFull = c.ID.Term
	}
	if c.RowType != "" {
		termFull = c.RowType
	}
	var idx int
	idxRes, err := strconv.Atoi(c.ID.Index)
	if err == nil {
		idx = idxRes
	}
	term := filepath.Base(termFull)
	term = strings.ToLower(term)
	coreData := CoreData{
		Index:      idx,
		Location:   c.Files.Location,
		TermFull:   termFull,
		Term:       term,
		FieldsData: make(map[string]FieldData),
		FieldsIdx:  make(map[int]FieldData),
	}

	for _, field := range c.Fields {
		idx = 0
		idxRes, err := strconv.Atoi(field.Index)

		if err == nil {
			idx = idxRes
		}

		term := filepath.Base(field.Term)
		term = strings.ToLower(term)
		fd := FieldData{
			Index:    idx,
			TermFull: field.Term,
			Term:     term,
		}
		coreData.FieldsData[term] = fd
		coreData.FieldsIdx[idx] = fd
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
		Location:   e.Files.Location,
		FieldsData: make(map[string]FieldData),
		FieldsIdx:  make(map[int]FieldData),
	}
	for _, field := range e.Fields {
		term := filepath.Base(field.Term)
		idx = 0
		idxRes, err := strconv.Atoi(field.Index)
		if err == nil {
			idx = idxRes
		}
		fd := FieldData{
			Index:    idx,
			TermFull: field.Term,
			Term:     term,
		}
		extData.FieldsData[term] = fd
		extData.FieldsIdx[idx] = fd
	}
	return extData
}

func stripExt(filename string) string {
	ext := len(filepath.Ext(filename))
	end := len(filename) - ext
	return filename[:end]
}
