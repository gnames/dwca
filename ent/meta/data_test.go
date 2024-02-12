package meta_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/dwca/ent/meta"
	"github.com/stretchr/testify/assert"
)

func TestToData(t *testing.T) {
	assert := assert.New(t)
	var m *meta.Meta
	var err error
	path := filepath.Join("..", "..", "testdata", "meta", "col.xml")
	f, err := os.Open(path)
	assert.Nil(err)
	defer f.Close()

	m, err = meta.New(f)
	assert.Nil(err)
	data := m.ToData()
	assert.NotNil(data)
	assert.Equal(0, data.Index)
	assert.Equal("", data.TermFull)
}
