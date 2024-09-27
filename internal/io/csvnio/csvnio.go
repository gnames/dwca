// package csvnio (CSV normal) uses Go's CSV library to read csv data.
package csvnio

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gnames/dwca/internal/ent"
	"github.com/gnames/dwca/internal/ent/dcfile"
	"github.com/gnames/gnfmt"
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
	if res.a.BadRowProcessing != gnfmt.ErrorBadRow {
		r.FieldsPerRecord = -1
	}
	res.r = r

	return res, nil
}

func NewWriter(attr ent.CSVAttr) (ent.CSVWriter, error) {
	res := &csvnio{a: attr}
	return res, nil
}

func (c *csvnio) ReadSlice(offset, limit int) ([][]string, error) {
	fieldsNum, lineNum, err := c.skipHeader()
	if err != nil {
		return nil, err
	}

	var res [][]string

	var count int
	for {
		lineNum++

		if limit > 0 && len(res) == limit {
			break
		}

		row, err := c.r.Read()
		if err == io.EOF {
			break
		}

		if fieldsNum == 0 {
			fieldsNum = len(row)
		}

		if err != nil {
			return nil, err
		}

		if offset > 0 && count <= offset {
			continue
		}
		rowFieldsNum := len(row)
		if fieldsNum == 0 {
			fieldsNum = rowFieldsNum
		}

		if rowFieldsNum != fieldsNum {
			skip := c.badRow(lineNum, fieldsNum, rowFieldsNum)
			if skip {
				continue
			} else {
				// set row to the required size
				row = gnfmt.NormRowSize(row, fieldsNum)
			}
		}

		count++
		res = append(res, row)
	}
	return res, nil
}

func (c *csvnio) skipHeader() (int, int, error) {
	var fieldsNum, lineNum int
	// ignore headers gif they are given
	if c.a.IgnoreHeader == "1" {
		lineNum++
		row, err := c.r.Read()
		if err != nil {
			return 0, 0, err
		}
		fieldsNum = len(row)
	}
	return fieldsNum, lineNum, nil
}

func (c *csvnio) badRow(
	lineNum, fieldsNum, rowFieldsNum int,
) bool {
	switch c.a.BadRowProcessing {
	case gnfmt.SkipBadRow:
		slog.Warn(
			"Wrong number of fields, SKIPPING row",
			"line", lineNum,
			"fieldsNum", fieldsNum,
			"rowFieldsNum", rowFieldsNum,
		)
		return true
	case gnfmt.ProcessBadRow:
		slog.Warn(
			"Wrong number of fields, PROCESSING the row anyway",
			"line", lineNum,
			"fieldsNum", fieldsNum,
			"rowFieldsNum", rowFieldsNum,
		)
	}
	return false
}

func (c *csvnio) Read(
	ctx context.Context,
	ch chan<- []string,
) (int, error) {
	// ignore headers if they are given
	fieldsNum, lineNum, err := c.skipHeader()
	if err != nil {
		return 0, err
	}

	var count int64
	for {
		lineNum++
		row, err := c.r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		rowFieldsNum := len(row)

		if fieldsNum == 0 {
			fieldsNum = rowFieldsNum
		}

		if fieldsNum != rowFieldsNum {
			skip := c.badRow(lineNum, fieldsNum, rowFieldsNum)
			if skip {
				continue
			} else {
				row = gnfmt.NormRowSize(row, fieldsNum)
			}
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
