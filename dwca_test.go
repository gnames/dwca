package dwca_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/dwca"
	"github.com/gnames/dwca/internal/ent/dcfile"
	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	assert := assert.New(t)
	res, err := dwca.Factory(filepath.Join("testdata", "data.zip"))
	assert.Nil(err)
	assert.Implements((*dwca.Archive)(nil), res)

	badPath := filepath.Join("testdata", "dont_exist.zip")
	res, err = dwca.Factory(badPath)
	assert.NotNil(err)
	assert.Nil(res)
	_, ok := err.(*dcfile.ErrFileNotFound)
	assert.True(ok)
}
