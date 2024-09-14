package factory

import (
	"github.com/gnames/dwca/internal/ent"
	"github.com/gnames/dwca/internal/io/csvnio"
	"github.com/gnames/dwca/internal/io/csvsio"
)

func CSVReader(attr ent.CSVAttr) (ent.CSVReader, error) {
	var res ent.CSVReader
	var err error
	if attr.Quote == "" {
		res, err = csvsio.New(attr)
	} else {
		res, err = csvnio.New(attr)
	}

	if err != nil {
		return nil, err
	}

	return res, err
}

func CSVWriter(attr ent.CSVAttr) (ent.CSVWriter, error) {
	res, err := csvnio.NewWriter(attr)
	if err != nil {
		return nil, err
	}
	return res, nil
}
