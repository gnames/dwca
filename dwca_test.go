package dwca_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/dwca"
	"github.com/gnames/dwca/ent/dcfile"
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
	assert.ErrorIs(err, dcfile.ErrFileNotFound{Path: badPath})
	assert.Nil(res)
}

func TestExtractZip(t *testing.T) {
	assert := assert.New(t)
	dwcaFile, err := dwca.Factory(filepath.Join("testdata", "data.zip"))
	assert.Nil(err)
	err = dwcaFile.Extract()
	assert.Nil(err)

	fakePath := filepath.Join("testdata", "fake.zip")
	dwcaFile, err = dwca.Factory(fakePath)
	assert.Nil(err)
	err = dwcaFile.Extract()
	assert.NotNil(err)
}

func TestExtractTar(t *testing.T) {
	assert := assert.New(t)
	dwcaFile, err := dwca.Factory(filepath.Join("testdata", "data.tar"))
	assert.Nil(err)
	err = dwcaFile.Extract()
	assert.Nil(err)

	fakePath := filepath.Join("testdata", "fake.tar")
	dwcaFile, err = dwca.Factory(fakePath)
	assert.Nil(err)
	err = dwcaFile.Extract()
	assert.NotNil(err)
	_, ok := err.(dcfile.ErrExtract)
	assert.True(ok)
}

func TestExtractTarGz(t *testing.T) {
	assert := assert.New(t)
	dwcaFile, err := dwca.Factory(filepath.Join("testdata", "data.tar.gz"))
	assert.Nil(err)
	err = dwcaFile.Extract()
	assert.Nil(err)

	fakePath := filepath.Join("testdata", "fake.tar.gz")
	dwcaFile, err = dwca.Factory(fakePath)
	assert.Nil(err)
	err = dwcaFile.Extract()
	assert.NotNil(err)
	_, ok := err.(dcfile.ErrExtract)
	assert.True(ok)
}

func TestExtractTarBz2(t *testing.T) {
	assert := assert.New(t)
	dwcaFile, err := dwca.Factory(filepath.Join("testdata", "data.tar.bz2"))
	assert.Nil(err)
	err = dwcaFile.Extract()
	assert.Nil(err)

	fakePath := filepath.Join("testdata", "fake.tar.bz2")
	dwcaFile, err = dwca.Factory(fakePath)
	assert.Nil(err)
	err = dwcaFile.Extract()
	assert.NotNil(err)
	_, ok := err.(dcfile.ErrExtract)
	assert.True(ok)
}

func TestExtractTarXz(t *testing.T) {
	assert := assert.New(t)
	dwcaFile, err := dwca.Factory(filepath.Join("testdata", "data.tar.xz"))
	assert.Nil(err)
	err = dwcaFile.Extract()
	assert.Nil(err)

	fakePath := filepath.Join("testdata", "fake.tar.xz")
	dwcaFile, err = dwca.Factory(fakePath)
	assert.Nil(err)
	err = dwcaFile.Extract()
	assert.NotNil(err)
	_, ok := err.(dcfile.ErrExtract)
	assert.True(ok)
}
