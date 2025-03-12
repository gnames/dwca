package dwca

import "encoding/xml"

type Meta struct {
	XMLName    xml.Name     `xml:"archive"`
	EMLFile    string       `xml:"metadata,attr"`
	Core       *Core        `xml:"core"`
	Extensions []*Extension `xml:"extension"`
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

// Core includes Attr and any core-specific fields (like ID).
type Core struct {
	ID ID `xml:"id"`
	*Attr
}

// Extension includes Attr and any extension-specific fields (like CoreID).
type Extension struct {
	CoreID CoreID `xml:"coreid"`
	*Attr
}

// Files holds the location of files.
type Files struct {
	// Location provides path to a file.
	Location string `xml:"location"`
}

// ID holds the fields for the Core ID.
type ID struct {
	Index string `xml:"index,attr"`
	Idx   int    `xml:"-"`
	Term  string `xml:"term,attr"`
}

// CoreID holds the fields for the CoreID data.
type CoreID struct {
	Index string `xml:"index,attr"`
	Idx   int    `xml:"-"`
}

// Field holds the fields of the data.
type Field struct {
	// Index is the verbatim index of the field.
	Index string `xml:"index,attr,omitempty"`

	// Idx is the integer version of Index.
	Idx int `xml:"-"`

	// Term is the URI of the term
	// (e.g., http://rs.tdwg.org/dwc/terms/scientificName).
	Term string `xml:"term,attr"`

	// Default provides a default value for this field across all rows.
	// This value corresponds to the given Term.
	Default string `xml:"default,attr,omitempty"`
}
