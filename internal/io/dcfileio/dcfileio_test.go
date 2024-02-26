package dcfileio_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gnames/dwca/internal/ent/dcfile"
	"github.com/gnames/dwca/internal/io/dcfileio"
	"github.com/gnames/dwca/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		msg  string
		file string
		err  bool
	}{
		{"zip", "data.zip", false},
		{"fake zip", "fake.zip", true},

		{"tar", "data.tar", false},
		{"fake tar", "fake.tar", true},

		{"tar.gz", "data.tar.gz", false},
		{"fake tar.gz", "fake.tar.gz", true},

		{"tar.bz2", "data.tar.bz2", false},
		{"fake tar.bz2", "fake.tar.bz2", true},

		{"tar.xz", "data.tar.xz", false},
		{"fake tar.xz", "fake.tar.xz", true},

		{"unknown", "unknown", true},
	}

	for _, v := range tests {
		path := filepath.Join("..", "..", "..", "pkg", "testdata", v.file)
		cfg := config.New()
		df, err := dcfileio.New(cfg, path)
		assert.Nil(err)

		// delete old archive dir
		err = df.Init()
		assert.Nil(err)

		// extract new archive dir
		err = df.Extract()
		assert.Equal(v.err, err != nil, v.msg)
		if err == nil {
			assert.Nil(err, v.msg)
			assert.NotNil(df, v.msg)
			continue
		}

		_, ok := err.(*dcfile.ErrExtract)
		if v.msg == "unknown" {
			_, ok = err.(*dcfile.ErrUnknownArchiveType)
		}
		assert.True(ok, v.msg)
		err = df.Close()
		assert.Nil(err)
	}
}

func TestMetaDir(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg  string
		file string
		path string
		err  bool
	}{
		{"zip", "data.zip", "IndexFungorumDWC", false},
		{"tar", "data.tar", "extract", false},
		{"tar.gz", "data.tar.gz", "extract", false},
		{"tar.bz2", "data.tar.bz2", "extract", false},
		{"tar.xz", "data.tar.xz", "extract", false},
		{"duplicate", "meta_dupl.tar.gz", "", true},
		{"absent", "meta_absent.tar.gz", "", true},
	}
	for _, v := range tests {
		cfg := config.New()
		df, err := dcfileio.New(
			cfg,
			filepath.Join("..", "..", "..", "pkg", "testdata", v.file),
		)
		assert.Nil(err)

		// delte old archive dir
		err = df.Init()
		assert.Nil(err)

		// extract new archive dir
		err = df.Extract()
		assert.Nil(err)

		// find archive dir
		arcDir, err := df.ArchiveDir()
		if err == nil {
			assert.True(strings.HasSuffix(arcDir, v.path), v.msg)
			continue
		}

		assert.Equal("", arcDir)
		ok := false

		if v.msg == "absent" {
			_, ok = err.(*dcfile.ErrMetaFileNotFound)
		}
		if v.msg == "duplicate" {
			_, ok = err.(*dcfile.ErrMultipleMetaFiles)
		}
		assert.True(ok)
		err = df.Close()
		assert.Nil(err)
	}
}
