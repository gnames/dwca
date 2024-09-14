// package csvsio (CSV Simple) uses simple split by field separator.
package csvsio

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gnames/dwca/internal/ent"
	"github.com/gnames/dwca/internal/ent/dcfile"
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

func (c *csvsio) ReadSlice(offset, limit int) ([][]string, error) {
	// ignore headers gif they are given
	if c.a.IgnoreHeader == "1" {
		c.r.Scan()
	}
	var res [][]string
	var fieldsCount int
	var count int
	for c.r.Scan() {
		count++

		if limit > 0 && len(res) == limit {
			break
		}

		if offset > 0 && count <= offset {
			continue
		}

		line := c.r.Text()
		sep := string(c.a.ColSep)
		row := strings.Split(line, sep)
		if fieldsCount == 0 {
			fieldsCount = len(row)
		}

		if !c.a.WithSloppyCSV && fieldsCount != len(row) {
			return nil, fmt.Errorf("wrong number of fieds: '%s'", line)
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
	// ignore headers if they are given
	if c.a.IgnoreHeader == "1" {
		c.r.Scan()
	}

	var fieldsCount int
	var count int64
	for c.r.Scan() {
		count++
		if count%100_000 == 0 {
			fmt.Printf("\r%s", strings.Repeat(" ", 50))
			fmt.Printf("\rProcessed %s lines", humanize.Comma(count))
		}

		line := c.r.Text()
		sep := string(c.a.ColSep)
		row := strings.Split(line, sep)
		if fieldsCount == 0 {
			fieldsCount = len(row)
		}

		if !c.a.WithSloppyCSV && fieldsCount != len(row) {
			return 0, fmt.Errorf("wrong number of fieds: '%s'", line)
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
