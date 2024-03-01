package meta_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/dwca/pkg/ent/meta"
	"github.com/stretchr/testify/assert"
)

func TestHeadersGBIF(t *testing.T) {
	var m *meta.Meta
	assert := assert.New(t)
	path := filepath.Join("..", "..", "testdata", "meta", "gbif.xml")
	f, err := os.Open(path)
	assert.Nil(err)
	m, err = meta.New(f)
	idx := m.Core.ID.Idx
	res := meta.Headers(idx, m.Core.Fields)
	assert.Equal(23, len(res))
}
