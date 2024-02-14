package dwca_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/dwca"
	"github.com/gnames/dwca/internal/ent/dcfile"
	"github.com/gnames/dwca/internal/ent/diagn"
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

func TestCoreData(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg    string
		file   string
		offset int
		limit  int
		len    int
		res00  string
	}{
		{"nolimit,nooffset", "data.tar.gz", 0, 0, 587, "leptogastrinae:tid:127"},
		{"limit", "data.tar.gz", 0, 10, 10, "leptogastrinae:tid:127"},
		{"offset", "data.tar.gz", 1, 0, 586, "leptogastrinae:tid:42"},
		{"offset,limit", "data.tar.gz", 1, 10, 10, "leptogastrinae:tid:42"},
	}
	for _, v := range tests {
		path := filepath.Join("testdata", v.file)
		arc, err := dwca.Factory(path)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load()
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		data, err := arc.CoreSlice(v.offset, v.limit)
		assert.Nil(err)
		assert.Equal(v.len, len(data))
		assert.Equal(v.res00, data[0][0])
	}
}

func TestCoreStream(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg  string
		file string
		len  int
	}{
		{"tar.gz", "data.tar.gz", 587},
	}
	for _, v := range tests {
		path := filepath.Join("testdata", v.file)
		arc, err := dwca.Factory(path)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load()
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		ch := make(chan []string)
		go func() {
			err := arc.CoreStream(ch)
			assert.Nil(err)
		}()

		var count int
		for range ch {
			count++
		}
		assert.Equal(v.len, count)
	}
}

func TestExtensionData(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg    string
		file   string
		index  int
		offset int
		limit  int
		len    int
		res00  string
	}{
		{"nolimit,nooffset", "data.tar.gz", 0, 0, 0, 1, "leptogastrinae:tid:42"},
		{"limit", "data.tar.gz", 0, 0, 10, 1, "leptogastrinae:tid:42"},
		{"offset", "data.tar.gz", 0, 1, 0, 0, ""},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", v.file)
		arc, err := dwca.Factory(path)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load()
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		data, err := arc.ExtensionSlice(0, v.offset, v.limit)
		assert.Nil(err)
		assert.Equal(v.len, len(data))
		if len(data) > 0 {
			assert.Equal(v.res00, data[0][0])
		}
	}
}

func TestExtensionStream(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg   string
		file  string
		index int
		len   int
	}{
		{"tar.gz", "data.tar.gz", 0, 1},
	}
	for _, v := range tests {
		path := filepath.Join("testdata", v.file)
		arc, err := dwca.Factory(path)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load()
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		ch := make(chan []string)
		go func() {
			err := arc.ExtensionStream(v.index, ch)
			assert.Nil(err)
		}()

		var count int
		for range ch {
			count++
		}
		assert.Equal(v.len, count)
	}
}

func TestSciNameDiagnose(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		msg    string
		file   string
		snType diagn.SciNameType
	}{
		{"can", "canonical.tar.gz", diagn.SciNameCanonical},
		{"full+auth", "full-auth-dup.tar.gz", diagn.SciNameFull},
		{"full", "unknown.tar.gz", diagn.SciNameUnknown},
		{"full", "composite.tar.gz", diagn.SciNameComposite},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", "diagn", "scinames", v.file)
		arc, err := dwca.Factory(path)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load()
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		diag, err := arc.Diagnose()
		assert.Nil(err)
		assert.Equal(v.snType.String(), diag.SciNameType.String())
	}
}

func TestSynDiagnose(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		msg    string
		file   string
		snType diagn.SynonymType
	}{
		{"ext", "in_extension.tar.gz", diagn.SynExtension},
		{"accepted", "in_core_accepted.tar.gz", diagn.SynAcceptedID},
		{"hierarchy", "hierarchy_deprecated.tar.gz", diagn.SynHierarchy},
		{"hierarchy", "hierarchy.tar.gz", diagn.SynHierarchy},
		{"unknown", "unknown.tar.gz", diagn.SynUnknown},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", "diagn", "synonyms", v.file)
		arc, err := dwca.Factory(path)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load()
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		diag, err := arc.Diagnose()
		assert.Nil(err)
		assert.Equal(v.snType.String(), diag.SynonymType.String())
	}
}

func TestHierDiagnose(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		msg   string
		file  string
		hType diagn.HierType
	}{
		{"tree", "tree.tar.gz", diagn.HierTree},
		{"tree", "tree_depr.tar.gz", diagn.HierTree},
		{"flat", "flat.tar.gz", diagn.HierFlat},
		{"unknown", "unknown.tar.gz", diagn.HierUnknown},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", "diagn", "hierarchy", v.file)
		arc, err := dwca.Factory(path)
		assert.Nil(err)
		assert.Implements((*dwca.Archive)(nil), arc)

		err = arc.Load()
		assert.Nil(err)

		meta := arc.Meta()
		assert.NotNil(meta)

		diag, err := arc.Diagnose()
		assert.Nil(err)
		assert.Equal(v.hType.String(), diag.HierType.String())
	}
}
