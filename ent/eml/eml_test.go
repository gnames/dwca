package eml_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/dwca/ent/eml"
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

func TestSmall(t *testing.T) {
	var e *eml.EML
	assert := assert.New(t)
	path := filepath.Join("..", "..", "testdata", "eml", "small.xml")
	f, err := os.Open(path)
	assert.Nil(err)
	defer f.Close()

	e, err = eml.New(f)
	assert.Nil(err)
	assert.Equal("123", e.Dataset.ID)
	assert.Contains(e.Dataset.Title, "Fungorum")
	assert.Contains(e.Dataset.Abstract.Para, "(previously ICBN)")
	creators := e.Dataset.Creators
	assert.Equal(1, len(creators))
	assert.Nil(creators[0].OrganizationName)
}

func TestMedium(t *testing.T) {
	var e *eml.EML
	assert := assert.New(t)
	path := filepath.Join("..", "..", "testdata", "eml", "medium.xml")
	f, err := os.Open(path)
	assert.Nil(err)
	defer f.Close()

	e, err = eml.New(f)
	assert.Nil(err)
	assert.Equal("http://dx.doi.org/10.14284/170", e.Dataset.ID)
	assert.Equal("10.14284/170", e.Dataset.AlternativeIdentifier.Value)

	creators := e.Dataset.Creators
	assert.Equal(1, len(creators))
	assert.Equal("", creators[0].Scope)
	assert.Contains(creators[0].OrganizationName.Value, "WoRMS")
	assert.Nil(creators[0].IndividualName)
	assert.Contains(creators[0].ElectronicMailAddress, "marinespecies.org")

	mdProvs := e.Dataset.MetadataProviders
	assert.Equal(1, len(mdProvs))
	assert.Nil(mdProvs[0].IndividualName)
	assert.Contains(mdProvs[0].OrganizationName.Value, "WoRMS")
	assert.Contains(mdProvs[0].ElectronicMailAddress, "marinespecies.org")
	assert.Contains(mdProvs[0].OnlineURL, "http")

	assParties := e.Dataset.AssociatedParties
	assert.Equal(15, len(assParties))
	ap1 := assParties[0]
	assert.NotNil(ap1.IndividualName)
	assert.NotNil(ap1.OrganizationName)
	assert.NotNil(ap1.Address)
	assert.Equal(1, len(ap1.Roles))

	assert.Contains(e.Dataset.PubDate, "2024")
	assert.Equal("eng", e.Dataset.Language)
	assert.Contains(e.Dataset.Abstract.Para, "An authoritative")

	kwSets := e.Dataset.KeywodSets
	assert.Equal(3, len(kwSets))
	kw1 := kwSets[0]
	assert.Equal(5, len(kw1.Keywords))
	assert.Equal("taxonomic", kw1.Keywords[0].Value)
	assert.Contains(e.Dataset.IntellectualRights.Para, "This work is licensed")

	cover := e.Dataset.Coverage
	assert.NotNil(cover)
	assert.NotNil(cover.GeographicCoverage)
	assert.NotNil(cover.TemporalCoverage)

	cont := e.Dataset.Contacts
	assert.Equal(1, len(cont))
	assert.Nil(cont[0].IndividualName)
	assert.Contains(cont[0].OrganizationName.Value, "World Register")
	assert.Contains(cont[0].ElectronicMailAddress, "marinespecies.org")
}

type badReader struct{}

func (b badReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("bad reader")
}
