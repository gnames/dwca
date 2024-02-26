package dwca_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/dwca/internal/ent/diagn"
	dwca "github.com/gnames/dwca/pkg"
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
		arc, err := dwca.Factory(path)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load()
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		err = arc.NormalizedDwCA(path)
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
		// {"tar", "flat.tar.gz", "flat.tar.gz"},
	}

	for _, v := range tests {
		ari := append([]string{"testdata", "diagn", "hierarchy"}, v.in)
		path := filepath.Join(ari...)
		arc, err := dwca.Factory(path)
		assert.Nil(err)
		err = arc.Load()
		assert.Nil(err)
		err = arc.NormalizedDwCA(path)
		assert.Nil(err)
		outPath := filepath.Join(arc.Config().DownloadPath, v.out)
		if v.msg == "zip" {
			err = arc.ZipNormalizedDwCA(outPath)
			assert.Nil(err)
		} else {
			err = arc.TarGzNormalizedDwCA(outPath)
			assert.Nil(err)
		}
	}
}
