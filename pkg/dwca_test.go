package dwca_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gnames/dwca/internal/ent/dcfile"
	"github.com/gnames/dwca/internal/ent/diagn"
	dwca "github.com/gnames/dwca/pkg"
	"github.com/gnames/dwca/pkg/config"
	"github.com/gnames/gnfmt"
	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	assert := assert.New(t)
	cfg := config.New()
	arc, err := dwca.Factory(filepath.Join("testdata", "data.zip"), cfg)
	assert.Nil(err)
	assert.Implements((*dwca.Archive)(nil), arc)

	badPath := filepath.Join("testdata", "dont_exist.zip")
	arc, err = dwca.Factory(badPath, cfg)
	assert.NotNil(err)
	assert.Nil(arc)
	_, ok := err.(*dcfile.ErrFileNotFound)
	assert.True(ok)
}

func TestStricterCSV(t *testing.T) {
	assert := assert.New(t)
	path := filepath.Join("testdata", "data.tar.gz")
	cfg := config.New()
	arc, err := dwca.Factory(path, cfg)
	assert.Nil(err)
	assert.Implements((*dwca.Archive)(nil), arc)

	err = arc.Load(cfg.ExtractPath)
	// breaks on diagnostics stage
	assert.NotNil(err)
}

func TestCoreData(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg    string
		file   string
		offset int
		limit  int
		len    int
		res00  string
	}{
		{"nolimit,nooffset", "data.tar.gz", 0, 0, 587, "leptogastrinae:tid:127"},
		{"pipe", "data_pipe.tar.gz", 0, 0, 587, "leptogastrinae:tid:127"},
		{"limit", "data.tar.gz", 0, 10, 10, "leptogastrinae:tid:127"},
		{"offset", "data.tar.gz", 1, 0, 586, "leptogastrinae:tid:42"},
		{"offset,limit", "data.tar.gz", 1, 10, 10, "leptogastrinae:tid:42"},
	}
	for _, v := range tests {
		path := filepath.Join("testdata", v.file)
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
		err = arc.Close()
		assert.Nil(err)
	}
}

// TestPreferTreeHierarchy checks if tree hierarchy is created if there is
// flat hierarchy also possible.
func TestPreferTreeHierarchy(t *testing.T) {
	assert := assert.New(t)
	path := filepath.Join("testdata", "vascan.zip")
	cfg := config.New()
	arc, err := dwca.Factory(path, cfg)
	assert.Nil(err)
	assert.Implements((*dwca.Archive)(nil), arc)

	err = arc.Load(cfg.ExtractPath)
	assert.Nil(err)

	err = arc.Normalize()
	assert.Nil(err)

	path = filepath.Join(cfg.OutputPath, "taxon.txt")
	bs, err := os.ReadFile(path)
	assert.Nil(err)

	rows := strings.Split(string(bs), "\n")
	assert.Less(1000, len(rows))
	// IDs look like `73|26|25|128|1142` and can only come from tree
	// hierachy.
	for _, v := range rows {
		if strings.Contains(v, "family|genus") {
			assert.Regexp(`\d+\|\d+\|\d+\|\d+`, v)
			break
		}
	}
	fmt.Println(path)

	err = arc.Close()
	assert.Nil(err)
}

// TestDomain checks if Domain is included in flat hierarchy.
func TestDomain(t *testing.T) {
	assert := assert.New(t)
	path := filepath.Join("testdata", "domain.tar.gz")
	cfg := config.New()
	arc, err := dwca.Factory(path, cfg)
	assert.Nil(err)
	assert.Implements((*dwca.Archive)(nil), arc)

	err = arc.Load(cfg.ExtractPath)
	assert.Nil(err)

	err = arc.Normalize()
	assert.Nil(err)

	path = filepath.Join(cfg.OutputPath, "taxa.txt")
	bs, err := os.ReadFile(path)
	assert.Nil(err)

	rows := strings.Split(string(bs), "\n")
	assert.Equal(8, len(rows))
	assert.Contains(rows[2], "domain")

	err = arc.Close()
	assert.Nil(err)
}

func TestCoreStream(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg  string
		file string
		len  int
	}{
		{"tar.gz", "data.tar.gz", 587},
		{"pipe", "data_pipe.tar.gz", 587},
	}
	for _, v := range tests {
		path := filepath.Join("testdata", v.file)
		cfg := config.New(config.OptWrongFieldsNum(gnfmt.ProcessBadRow))
		arc, err := dwca.Factory(path, cfg)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load(cfg.ExtractPath)
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		ch := make(chan []string)
		go func() {
			_, err = arc.CoreStream(context.Background(), ch)
			assert.Nil(err)
		}()

		var count int
		for range ch {
			count++
		}
		assert.Equal(v.len, count)
		err = arc.Close()
		assert.Nil(err)
	}
}

func TestExtensionData(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg    string
		file   string
		index  int
		offset int
		limit  int
		len    int
		res00  string
	}{
		{"nolimit,nooffset", "data.tar.gz", 0, 0, 0, 1, "leptogastrinae:tid:42"},
		{"pipe", "data_pipe.tar.gz", 0, 0, 0, 1, "leptogastrinae:tid:42"},
		{"limit", "data.tar.gz", 0, 0, 10, 1, "leptogastrinae:tid:42"},
		{"offset", "data.tar.gz", 0, 1, 0, 0, ""},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", v.file)
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
		if len(data) > 0 {
			assert.Equal(v.res00, data[0][0])
		}
		err = arc.Close()
		assert.Nil(err)
	}
}

func TestExtensionStream(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg   string
		file  string
		index int
		len   int
	}{
		{"tar.gz", "data.tar.gz", 0, 1},
		{"pipe", "data_pipe.tar.gz", 0, 1},
	}
	ctx := context.Background()
	for _, v := range tests {
		path := filepath.Join("testdata", v.file)
		cfg := config.New(config.OptWrongFieldsNum(gnfmt.ProcessBadRow))
		arc, err := dwca.Factory(path, cfg)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load(cfg.ExtractPath)
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		ch := make(chan []string)
		go func() {
			_, err = arc.ExtensionStream(ctx, v.index, ch)
			assert.Nil(err)
		}()

		var count int
		for range ch {
			count++
		}
		assert.Equal(v.len, count)
		err = arc.Close()
		assert.Nil(err)
	}
}

func TestSynDiagnose(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		msg    string
		file   string
		snType diagn.SynonymType
	}{
		{"ext", "in_extension.tar.gz", diagn.SynExtension},
		{"accepted", "in_core_accepted.tar.gz", diagn.SynAcceptedID},
		{"hierarchy", "hierarchy_deprecated.tar.gz", diagn.SynHierarchy},
		{"hierarchy", "hierarchy.tar.gz", diagn.SynHierarchy},
		{"unknown", "unknown.tar.gz", diagn.SynUnknown},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", "diagn", "synonyms", v.file)
		cfg := config.New(config.OptWrongFieldsNum(gnfmt.ProcessBadRow))
		arc, err := dwca.Factory(path, cfg)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load(cfg.ExtractPath)
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		err = arc.Close()
		assert.Nil(err)
	}
}
