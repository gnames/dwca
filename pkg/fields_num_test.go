package dwca_test

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"
	"sync"
	"testing"

	dwca "github.com/gnames/dwca/pkg"
	"github.com/gnames/dwca/pkg/config"
	"github.com/gnames/gnfmt"
	"github.com/stretchr/testify/assert"
)

func TestNormCoreData(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	assert := assert.New(t)
	tests := []struct {
		msg    string
		file   string
		offset int
		limit  int
		len    int
		res00  string
	}{
		{"csv", "csv-norm.tar.gz", 0, 0, 10, "2"},
		{"tsv", "tsv-norm.tar.gz", 0, 0, 10, "2"},
		{"err-more-csv", "csv-more.tar.gz", 0, 0, 0, ""},
		{"err-less-csv", "csv-less.tar.gz", 0, 0, 0, ""},
		{"err-more-tsv", "tsv-more.tar.gz", 0, 0, 0, ""},
		{"err-less-tsv", "tsv-less.tar.gz", 0, 0, 0, ""},
	}
	for _, v := range tests {
		path := filepath.Join("testdata", "fldnum", v.file)
		cfg := config.New()
		arc, err := dwca.Factory(path, cfg)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load(cfg.ExtractPath)
		if v.res00 == "" {
			assert.NotNil(err)
			continue
		}
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		data, err := arc.CoreSlice(v.offset, v.limit)
		assert.Nil(err)
		assert.Equal(v.len, len(data))
		assert.Equal(v.res00, data[0][0])

		chIn := make(chan []string)
		var wg sync.WaitGroup
		wg.Add(1)

		data = nil
		go func() {
			defer wg.Done()
			for row := range chIn {
				data = append(data, row)
			}
		}()
		var count int
		count, err = arc.CoreStream(context.Background(), chIn)
		assert.Nil(err)
		wg.Wait()
		assert.Equal(v.len, len(data))
		assert.Equal(v.len, count)
		assert.Equal(v.res00, data[0][0])
		err = arc.Close()
		assert.Nil(err)
	}
}

func TestNormExtData(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	assert := assert.New(t)
	tests := []struct {
		msg    string
		file   string
		offset int
		limit  int
		len    int
		res00  string
	}{
		{"csv", "csv-norm.tar.gz", 0, 0, 20, "2"},
		{"tsv", "tsv-norm.tar.gz", 0, 0, 20, "2"},
		{"err-more-csv", "csv-more.tar.gz", 0, 0, 0, ""},
		{"err-less-csv", "csv-less.tar.gz", 0, 0, 0, ""},
		{"err-more-tsv", "tsv-more.tar.gz", 0, 0, 0, ""},
		{"err-less-tsv", "tsv-less.tar.gz", 0, 0, 0, ""},
	}
	for _, v := range tests {
		path := filepath.Join("testdata", "fldnum", v.file)
		cfg := config.New()
		arc, err := dwca.Factory(path, cfg)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load(cfg.ExtractPath)
		if v.res00 == "" {
			assert.NotNil(err)
			continue
		}
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		data, err := arc.ExtensionSlice(0, v.offset, v.limit)
		assert.Nil(err)
		assert.Equal(v.len, len(data))
		assert.Equal(v.res00, data[0][0])

		chIn := make(chan []string)
		var wg sync.WaitGroup
		wg.Add(1)

		data = nil
		go func() {
			defer wg.Done()
			for row := range chIn {
				data = append(data, row)
			}
		}()
		var count int
		count, err = arc.ExtensionStream(context.Background(), 0, chIn)
		assert.Nil(err)
		wg.Wait()
		assert.Equal(v.len, len(data))
		assert.Equal(v.len, count)
		assert.Equal(v.res00, data[0][0])
		err = arc.Close()
		assert.Nil(err)
	}
}

func TestSkipCoreRows(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	assert := assert.New(t)
	tests := []struct {
		msg    string
		file   string
		offset int
		limit  int
		len    int
		res00  string
	}{
		{"csv", "csv-norm.tar.gz", 0, 0, 10, "2"},
		{"tsv", "tsv-norm.tar.gz", 0, 0, 10, "2"},
		{"more-csv", "csv-more.tar.gz", 0, 0, 9, "2"},
		{"less-csv", "csv-less.tar.gz", 0, 0, 9, "2"},
		{"more-tsv", "tsv-more.tar.gz", 0, 0, 9, "2"},
		{"less-tsv", "tsv-less.tar.gz", 0, 0, 9, "2"},
	}
	for _, v := range tests {
		path := filepath.Join("testdata", "fldnum", v.file)
		cfg := config.New(config.OptWrongFieldsNum(gnfmt.SkipBadRow))
		arc, err := dwca.Factory(path, cfg)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load(cfg.ExtractPath)
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		data, err := arc.CoreSlice(v.offset, v.limit)
		assert.Nil(err)
		assert.Equal(v.len, len(data))
		assert.Equal(v.res00, data[0][0])

		chIn := make(chan []string)
		var wg sync.WaitGroup
		wg.Add(1)

		data = nil
		go func() {
			defer wg.Done()
			for row := range chIn {
				data = append(data, row)
			}
		}()
		var count int
		count, err = arc.CoreStream(context.Background(), chIn)
		assert.Nil(err)
		wg.Wait()
		assert.Equal(v.len, len(data))
		assert.Equal(v.len, count)
		assert.Equal(v.res00, data[0][0])
		err = arc.Close()
		assert.Nil(err)
	}
}

func TestSkipExtRows(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	assert := assert.New(t)
	tests := []struct {
		msg    string
		file   string
		offset int
		limit  int
		len    int
		res00  string
	}{
		{"csv", "csv-norm.tar.gz", 0, 0, 20, "2"},
		{"tsv", "tsv-norm.tar.gz", 0, 0, 20, "2"},
		{"more-csv", "csv-more.tar.gz", 0, 0, 18, "2"},
		{"less-csv", "csv-less.tar.gz", 0, 0, 18, "2"},
		{"more-tsv", "tsv-more.tar.gz", 0, 0, 18, "2"},
		{"less-tsv", "tsv-less.tar.gz", 0, 0, 18, "2"},
	}
	for _, v := range tests {
		path := filepath.Join("testdata", "fldnum", v.file)
		cfg := config.New(config.OptWrongFieldsNum(gnfmt.SkipBadRow))
		arc, err := dwca.Factory(path, cfg)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load(cfg.ExtractPath)
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		data, err := arc.ExtensionSlice(0, v.offset, v.limit)
		assert.Nil(err)
		assert.Equal(v.len, len(data))
		assert.Equal(v.res00, data[0][0])

		chIn := make(chan []string)
		var wg sync.WaitGroup
		wg.Add(1)

		data = nil
		go func() {
			defer wg.Done()
			for row := range chIn {
				data = append(data, row)
			}
		}()
		var count int
		count, err = arc.ExtensionStream(context.Background(), 0, chIn)
		assert.Nil(err)
		wg.Wait()
		assert.Equal(v.len, len(data))
		assert.Equal(v.len, count)
		assert.Equal(v.res00, data[0][0])
		err = arc.Close()
		assert.Nil(err)
	}
}

func TestProcessCoreRows(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg    string
		file   string
		offset int
		limit  int
		len    int
		res00  string
	}{
		{"csv", "csv-norm.tar.gz", 0, 0, 10, "2"},
		{"tsv", "tsv-norm.tar.gz", 0, 0, 10, "2"},
		{"more-csv", "csv-more.tar.gz", 0, 0, 10, "2"},
		{"less-csv", "csv-less.tar.gz", 0, 0, 10, "2"},
		{"more-tsv", "tsv-more.tar.gz", 0, 0, 10, "2"},
		{"less-tsv", "tsv-less.tar.gz", 0, 0, 10, "2"},
	}
	for _, v := range tests {
		path := filepath.Join("testdata", "fldnum", v.file)
		cfg := config.New(config.OptWrongFieldsNum(gnfmt.ProcessBadRow))
		arc, err := dwca.Factory(path, cfg)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load(cfg.ExtractPath)
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		data, err := arc.CoreSlice(v.offset, v.limit)
		assert.Nil(err)
		assert.Equal(v.len, len(data))
		assert.Equal(v.res00, data[0][0])

		chIn := make(chan []string)
		var wg sync.WaitGroup
		wg.Add(1)

		data = nil
		go func() {
			defer wg.Done()
			for row := range chIn {
				data = append(data, row)
			}
		}()
		var count int
		count, err = arc.CoreStream(context.Background(), chIn)
		assert.Nil(err)
		wg.Wait()
		assert.Equal(v.len, len(data))
		assert.Equal(v.len, count)
		assert.Equal(v.res00, data[0][0])

		err = arc.Close()
		assert.Nil(err)
	}
}

func TestProcessExtRows(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg    string
		file   string
		offset int
		limit  int
		len    int
		res00  string
	}{
		{"csv", "csv-norm.tar.gz", 0, 0, 20, "2"},
		{"tsv", "tsv-norm.tar.gz", 0, 0, 20, "2"},
		{"more-csv", "csv-more.tar.gz", 0, 0, 20, "2"},
		{"less-csv", "csv-less.tar.gz", 0, 0, 20, "2"},
		{"more-tsv", "tsv-more.tar.gz", 0, 0, 20, "2"},
		{"less-tsv", "tsv-less.tar.gz", 0, 0, 20, "2"},
	}
	for _, v := range tests {
		path := filepath.Join("testdata", "fldnum", v.file)
		cfg := config.New(config.OptWrongFieldsNum(gnfmt.ProcessBadRow))
		arc, err := dwca.Factory(path, cfg)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load(cfg.ExtractPath)
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		data, err := arc.ExtensionSlice(0, v.offset, v.limit)
		assert.Nil(err)
		assert.Equal(v.len, len(data))
		assert.Equal(v.res00, data[0][0])

		chIn := make(chan []string)
		var wg sync.WaitGroup
		wg.Add(1)

		data = nil
		go func() {
			defer wg.Done()
			for row := range chIn {
				data = append(data, row)
			}
		}()
		var count int
		count, err = arc.ExtensionStream(context.Background(), 0, chIn)
		assert.Nil(err)
		wg.Wait()
		assert.Equal(v.len, len(data))
		assert.Equal(v.len, count)
		assert.Equal(v.res00, data[0][0])

		err = arc.Close()
		assert.Nil(err)
	}
}
