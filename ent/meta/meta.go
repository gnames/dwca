package meta

import "encoding/xml"

type Meta struct {
	Archive `xml:"archive"`
}

type Archive struct {
	XMLName    xml.Name    `xml:"archive"`
	EMLFile    string      `xml:"metadata,attr"`
	Core       Core        `xml:"core"`
	Extensions []Extension `xml:"extension"`
}

// Attr holds the common fields for Core and Extension.
type Attr struct {
	Encoding           string  `xml:"encoding,attr"`
	FieldsTerminatedBy string  `xml:"fieldsTerminatedBy,attr"`
	LinesTerminatedBy  string  `xml:"linesTerminatedBy,attr"`
	FieldsEnclosedBy   string  `xml:"fieldsEnclosedBy,attr"`
	IgnoreHeaderLines  string  `xml:"ignoreHeaderLines,attr"`
	RowType            string  `xml:"rowType,attr"`
	Files              Files   `xml:"files"`
	Fields             []Field `xml:"field"`
}

// Core includes CommonElement and any core-specific fields (like ID).
type Core struct {
	ID ID `xml:"id"`
	Attr
}

// Extension includes CommonElement and any extension-specific fields (like CoreID).
type Extension struct {
	Attr
	CoreID CoreID `xml:"coreid"`
}

type Files struct {
	Location string `xml:"location"`
}

type ID struct {
	Index string `xml:"index,attr"`
	Term  string `xml:"term,attr"`
}

type CoreID struct {
	Index string `xml:"index,attr"`
}

type Field struct {
	Index string `xml:"index,attr"`
	Term  string `xml:"term,attr"`
}
