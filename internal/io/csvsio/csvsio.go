// package csvsio (CSV Simple) uses simple split by field separator.
package csvsio

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gnames/dwca/internal/ent"
	"github.com/gnames/dwca/internal/ent/dcfile"
	"github.com/gnames/gnfmt"
)

type csvsio struct {
	a ent.CSVAttr
	f *os.File
	r *bufio.Scanner
	w *bufio.Writer
}

func New(attr ent.CSVAttr) (ent.CSVReader, error) {
	res := &csvsio{a: attr}

	f, err := os.Open(res.a.Path)
	if err != nil {
		return nil, err
	}
	res.f = f

	res.r = bufio.NewScanner(f)

	return res, nil
}

func NewWriter(attr ent.CSVAttr) (ent.CSVWriter, error) {
	res := &csvsio{a: attr}
	f, err := os.Open(res.a.Path)
	if err != nil {
		return nil, err
	}
	res.f = f

	res.w = bufio.NewWriter(f)
	return res, nil
}

func (c *csvsio) skipHeader() (int, int) {
	var fieldsNum, lineNum int
	// ignore headers gif they are given
	if c.a.IgnoreHeader == "1" {
		lineNum++
		c.r.Scan()
		line := c.r.Text()
		sep := string(c.a.ColSep)
		row := strings.Split(line, sep)
		fieldsNum = len(row)
	}
	return fieldsNum, lineNum
}

func (c *csvsio) badRow(
	lineNum, fieldsNum, rowFieldsNum int,
) (bool, error) {
	switch c.a.BadRowProcessing {
	case gnfmt.ErrorBadRow:
		err := fmt.Errorf("wrong number of fieds: '%d'", lineNum)
		slog.Error("Bad row",
			"line", lineNum,
			"fieldsNum", fieldsNum,
			"rowFieldsNum", rowFieldsNum,
			"error", err,
		)
		return false, err
	case gnfmt.SkipBadRow:
		slog.Warn(
			"Wrong number of fields, SKIPPING row",
			"line", lineNum,
			"fieldsNum", fieldsNum,
			"rowFieldsNum", rowFieldsNum,
		)
		return true, nil
	case gnfmt.ProcessBadRow:
		slog.Warn(
			"Wrong number of fields, PROCESSING the row anyway",
			"line", lineNum,
			"fieldsNum", fieldsNum,
			"rowFieldsNum", rowFieldsNum,
		)
	}
	return false, nil
}

func (c *csvsio) ReadSlice(offset, limit int) ([][]string, error) {

	fieldsNum, lineNum := c.skipHeader()

	var res [][]string
	var count int
	for c.r.Scan() {
		count++
		lineNum++

		if limit > 0 && len(res) == limit {
			break
		}

		if offset > 0 && count <= offset {
			continue
		}

		line := c.r.Text()
		sep := string(c.a.ColSep)
		row := strings.Split(line, sep)
		rowFieldsNum := len(row)
		if fieldsNum == 0 {
			fieldsNum = rowFieldsNum
		}

		if fieldsNum != rowFieldsNum {
			skip, err := c.badRow(lineNum, fieldsNum, rowFieldsNum)
			if skip {
				continue
			}
			if err != nil {
				return nil, err
			}
			row = gnfmt.NormRowSize(row, fieldsNum)
		}

		res = append(res, row)
	}

	if err := c.r.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *csvsio) Read(
	ctx context.Context,
	ch chan<- []string,
) (int, error) {
	fieldsNum, lineNum := c.skipHeader()

	var count int64
	for c.r.Scan() {
		if count%100_000 == 0 {
			fmt.Printf("\r%s", strings.Repeat(" ", 50))
			fmt.Printf("\rProcessed %s lines", humanize.Comma(count))
		}

		line := c.r.Text()
		sep := string(c.a.ColSep)
		row := strings.Split(line, sep)
		rowFieldsNum := len(row)
		if fieldsNum == 0 {
			fieldsNum = rowFieldsNum
		}

		if fieldsNum != rowFieldsNum {
			skip, err := c.badRow(lineNum, fieldsNum, rowFieldsNum)
			if skip {
				continue
			}
			if err != nil {
				return 0, err
			}
			row = gnfmt.NormRowSize(row, fieldsNum)
		}

		select {
		case <-ctx.Done():
			return 0, &dcfile.ErrContext{Err: ctx.Err()}
		default:
			count++
			ch <- row
		}
	}

	fmt.Printf("\r%s\r", strings.Repeat(" ", 50))
	return int(count), nil
}

func (c *csvsio) Close() error {
	return c.f.Close()
}

func (c *csvsio) Write(
	ctx context.Context,
	outChan <-chan []string,
) error {
	f, err := os.Create(c.a.Path)
	if err != nil {
		return &dcfile.ErrSaveCSV{Err: err}
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	headers := strings.Join(c.a.Headers, string(c.a.ColSep)) + "\n"
	_, err = w.Write([]byte(headers))
	if err != nil {
		return &dcfile.ErrSaveCSV{Err: err}
	}
	for row := range outChan {
		line := strings.Join(row, string(c.a.ColSep)) + "\n"
		_, err = w.Write([]byte(line))
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
