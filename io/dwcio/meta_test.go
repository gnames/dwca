package dwcio_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/dwca"
	"github.com/stretchr/testify/assert"
)

func TestMetaCoL(t *testing.T) {
	assert := assert.New(t)
	var m *dwca.Meta
	_ = m
	var err error
	path := filepath.Join("..", "..", "testdata", "meta", "col.xml")
	f, err := os.Open(path)
	assert.Nil(err)
	defer f.Close()

	// m, err = meta.New(f)
	// assert.Nil(err)
	// assert.IsType(m, &meta.Meta{})
	// assert.Equal(m.EMLFile, "eml.xml")
	// assert.Equal(m.Core.Encoding, "utf-8")
	// assert.Equal(m.Core.Files.Location, "Taxon.tsv")
	// assert.Equal(m.Core.FieldsTerminatedBy, "\\t")
	// assert.Equal(m.Core.ID.Index, "0")
	// assert.Equal(m.Archive.Core.Fields[7].Index, "7")
	// assert.Equal("http://rs.tdwg.org/dwc/terms/taxonRank", m.Archive.Core.Fields[7].Term)
	// assert.Equal(3, len(m.Extensions))
	// assert.Equal(m.Extensions[0].CoreID.Index, "0")
	// assert.Equal("http://rs.gbif.org/terms/1.0/isExtinct", m.Extensions[1].Fields[1].Term)
}
