package dwca_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/dwca/internal/ent/diagn"
	dwca "github.com/gnames/dwca/pkg"
	"github.com/gnames/dwca/pkg/config"
	"github.com/gnames/gnfmt"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeDwCA(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		msg   string
		path  []string
		hType diagn.HierType
	}{
		{"tree", []string{"data.tar.gz"}, diagn.HierTree},
		{"flat", []string{"diagn", "hierarchy", "flat.tar.gz"}, diagn.HierFlat},
		{"myriatrix", []string{"myriatrix.tar.gz"}, diagn.HierTree},
	}

	for _, v := range tests {
		v.path = append([]string{"testdata"}, v.path...)
		path := filepath.Join(v.path...)
		cfg := config.New(config.OptWrongFieldsNum(gnfmt.ProcessBadRow))
		arc, err := dwca.Factory(path, cfg)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load(cfg.ExtractPath)
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)
		err = arc.Normalize()
		assert.Nil(err)
	}
}

func TestCompress(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg string
		in  string
		out string
	}{
		{"zip", "flat.tar.gz", "flat.zip"},
		{"tar", "flat.tar.gz", "flat.tar.gz"},
	}

	for _, v := range tests {
		ari := append([]string{"testdata", "diagn", "hierarchy"}, v.in)
		path := filepath.Join(ari...)
		cfg := config.New(config.OptWrongFieldsNum(gnfmt.ProcessBadRow))
		arc, err := dwca.Factory(path, cfg)
		assert.Nil(err)
		err = arc.Load(cfg.ExtractPath)
		assert.Nil(err)
		err = arc.Normalize()
		assert.Nil(err)
		outPath := filepath.Join(arc.Config().DownloadPath, v.out)
		if v.msg == "zip" {
			err = arc.ZipNormalized(outPath)
			assert.Nil(err)
		} else {
			err = arc.TarGzNormalized(outPath)
			assert.Nil(err)
		}
	}
}

func TestIndexNoField(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, file string
		fieldNum  int
	}{
		{"idx norm", "tree.tar.gz", 8},
		{"idx empty", "tree_no_index_info.tar.gz", 8},
		{"idx empty", "gbif-small.tar.gz", 24},
	}
	for _, v := range tests {
		path := filepath.Join("testdata", "diagn", "hierarchy", v.file)
		cfg := config.New(config.OptWrongFieldsNum(gnfmt.ProcessBadRow))
		arc, err := dwca.Factory(path, cfg)
		assert.Nil(err, v.msg)

		err = arc.Load(cfg.ExtractPath)
		assert.Nil(err, v.msg)

		err = arc.Normalize()
		assert.Nil(err, v.msg)

		arc, err = dwca.Factory("", cfg)
		assert.Nil(err, v.msg)

		err = arc.Load(cfg.OutputPath)
		assert.Nil(err, v.msg)
		var ary [][]string
		ary, err = arc.CoreSlice(0, 10)
		assert.Nil(err, v.msg)
		for _, fld := range ary {
			assert.Equal(v.fieldNum, len(fld), v.msg)
		}
	}
}

func TestAosBirds(t *testing.T) {
	assert := assert.New(t)
	path := filepath.Join("testdata", "aos-birds.tar.gz")
	cfg := config.New()
	arc, err := dwca.Factory(path, cfg)
	assert.Nil(err)

	err = arc.Load(cfg.ExtractPath)
	assert.Nil(err)

	err = arc.Normalize()
	assert.Nil(err)

	arc, err = dwca.FactoryOutput(cfg)
	assert.Nil(err)

	err = arc.Load(cfg.OutputPath)
	assert.Nil(err)
	var ary [][]string
	ary, err = arc.CoreSlice(0, 10)
	assert.Nil(err)
	assert.Equal(10, len(ary))
}
