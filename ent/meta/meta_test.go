package meta_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/dwca/ent/eml"
	"github.com/gnames/dwca/ent/meta"
	"github.com/stretchr/testify/assert"
)

func TestConstructorErrors(t *testing.T) {
	assert := assert.New(t)
	e, err := eml.New(badReader{})
	assert.NotNil(err)
	assert.IsType(&eml.ErrReader{}, err)
	assert.Nil(e)

	e, err = eml.New(bytes.NewReader([]byte("")))
	assert.NotNil(err)
	assert.IsType(&eml.ErrDecoder{}, err)
	assert.Nil(e)
}

// TestMetaCoL tests reading of a Catalogue of Life meta.xml file.
func TestMetaCoL(t *testing.T) {
	assert := assert.New(t)
	var m *meta.Meta
	var err error
	path := filepath.Join("..", "..", "testdata", "meta", "col.xml")
	f, err := os.Open(path)
	assert.Nil(err)
	defer f.Close()

	m, err = meta.New(f)
	assert.Nil(err)
	assert.IsType(m, &meta.Meta{})
	assert.Equal(m.EMLFile, "eml.xml")
	assert.Equal(m.Core.Encoding, "utf-8")
	assert.Equal(m.Core.Files.Location, "Taxon.tsv")
	assert.Equal(m.Core.FieldsTerminatedBy, "\\t")
	assert.Equal(m.Core.ID.Index, "0")
	assert.Equal(m.Archive.Core.Fields[7].Index, "7")
	assert.Equal("http://rs.tdwg.org/dwc/terms/taxonRank", m.Archive.Core.Fields[7].Term)
	assert.Equal(3, len(m.Extensions))
	assert.Equal(m.Extensions[0].CoreID.Index, "0")
	assert.Equal("http://rs.gbif.org/terms/1.0/isExtinct", m.Extensions[1].Fields[1].Term)
}

type badReader struct{}

func (b badReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("bad reader")
}
