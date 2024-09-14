package ent

import "context"

type CSVAttr struct {
	Headers       []string
	Path          string
	ColSep        rune
	Quote         string
	IgnoreHeader  string
	WithSloppyCSV bool
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
