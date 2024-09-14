// package csvnio (CSV normal) uses Go's CSV library to read csv data.
package csvnio

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gnames/dwca/internal/ent"
	"github.com/gnames/dwca/internal/ent/dcfile"
)

type csvnio struct {
	a ent.CSVAttr
	f *os.File
	r *csv.Reader
}

func New(attr ent.CSVAttr) (ent.CSVReader, error) {
	res := &csvnio{a: attr}

	f, err := os.Open(res.a.Path)
	if err != nil {
		return nil, err
	}
	res.f = f

	r := csv.NewReader(f)
	r.Comma = res.a.ColSep

	// allow variable number of fields
	r.FieldsPerRecord = -1
	res.r = r

	return res, nil
}

func NewWriter(attr ent.CSVAttr) (ent.CSVWriter, error) {
	res := &csvnio{a: attr}
	return res, nil
}

func (c *csvnio) ReadSlice(offset, limit int) ([][]string, error) {
	// ignore headers gif they are given
	if c.a.IgnoreHeader == "1" {
		c.r.Read()
	}
	var res [][]string

	var count int
	for {
		count++

		if limit > 0 && len(res) == limit {
			break
		}

		row, err := c.r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, &dcfile.ErrCoreRead{Err: err}
		}

		if offset > 0 && count <= offset {
			continue
		}
		res = append(res, row)
	}

	return res, nil
}

func (c *csvnio) Read(
	ctx context.Context,
	ch chan<- []string,
) (int, error) {
	// ignore headers if they are given
	if c.a.IgnoreHeader == "1" {
		c.r.Read()
	}

	var count int64
	for {
		row, err := c.r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, &dcfile.ErrExtensionRead{Err: err}
		}

		count++
		if count%100_000 == 0 {
			fmt.Printf("\r%s", strings.Repeat(" ", 50))
			fmt.Printf("\rProcessed %s lines", humanize.Comma(count))
		}

		select {
		case <-ctx.Done():
			return 0, &dcfile.ErrContext{Err: ctx.Err()}
		default:
			ch <- row
		}
	}

	fmt.Printf("\r%s\r", strings.Repeat(" ", 50))
	return int(count), nil
}

func (c *csvnio) Close() error {
	return c.f.Close()
}

func (c *csvnio) Write(
	ctx context.Context,
	outChan <-chan []string,
) error {
	f, err := os.Create(c.a.Path)
	if err != nil {
		return &dcfile.ErrSaveCSV{Err: err}
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = c.a.ColSep

	err = w.Write(c.a.Headers)
	if err != nil {
		return &dcfile.ErrSaveCSV{Err: err}
	}
	for row := range outChan {
		err = w.Write(row)
		if err != nil {
			for range outChan {
			}
			return &dcfile.ErrSaveCSV{Err: err}
		}
		select {
		case <-ctx.Done():
			for range outChan {
			}
			return &dcfile.ErrContext{Err: ctx.Err()}
		default:
		}
	}
	w.Flush()
	return nil
}
