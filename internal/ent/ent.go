package ent

import (
	"context"
)

// BadRow type describes different scenarios of processing rows with wrong
// number of fields.
type BadRow int

const (
	// ProcessBadRow means processing bad row hoping for the best.
	ProcessBadRow BadRow = iota

	// SkipBadRow means that rows with wrong number of fields will not be
	// processed.
	SkipBadRow

	// ErrorBadRow means that an error will be returned when a row with wrong
	// number of fields is encountered.
	ErrorBadRow
)

// CSVAttr describes a variety of configuration attributes for reading and
// writing CSV files.
type CSVAttr struct {
	// Headers contains names of fields to be placed to CSV file during
	// creation of DwC Archive.
	Headers []string

	// Path is the path to the CSV file.
	Path string

	// ColSep is the UTF-8 character used to separate fields from each other.
	ColSep rune

	// Quote (usually `"`) that escapes ColSep characters withing the fields.
	Quote string

	// IgnoreHeader indicates if there is a header row in the CSV file.
	// If header exists, its values will be ignored.
	IgnoreHeader string

	// BadRowProcessing determines a method for dealing with rows that have
	// wrong number of elements. The 'bad rows' would either be processed,
	// ignored, of break the execution of the program. Default is to raise an
	// error.
	BadRowProcessing BadRow
}

type CSVReader interface {
	ReadSlice(offset, limit int) ([][]string, error)
	Read(context.Context, chan<- []string) (int, error)
	Close() error
}

type CSVWriter interface {
	Write(ctx context.Context, ch <-chan []string) error
	Close() error
}
