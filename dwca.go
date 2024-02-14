package dwca

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gnames/dwca/config"
	"github.com/gnames/dwca/ent/eml"
	"github.com/gnames/dwca/ent/meta"
	"github.com/gnames/dwca/internal/ent/dcfile"
	"github.com/gnames/dwca/internal/ent/diagn"
	"github.com/gnames/gnparser"
)

type arch struct {
	cfg      config.Config
	dcFile   dcfile.DCFile
	metaData *meta.Meta
	emlData  *eml.EML
	diagn.Diagnostics
	gnpPool chan gnparser.GNparser
}

func New(cfg config.Config, df dcfile.DCFile) Archive {
	res := &arch{cfg: cfg, dcFile: df}
	poolSize := 5
	gnpPool := make(chan gnparser.GNparser, poolSize)
	for i := 0; i < poolSize; i++ {
		cfgGNP := gnparser.NewConfig()
		gnpPool <- gnparser.New(cfgGNP)
	}
	res.gnpPool = gnpPool
	return res
}

// Config returns the configuration object of the archive.
func (a *arch) Config() config.Config {
	return a.cfg
}

// Load extracts the archive and loads data for EML and Meta.
func (a *arch) Load() error {
	err := a.dcFile.Extract()
	if err != nil {
		return err
	}
	path, err := a.dcFile.ArchiveDir()
	if err != nil {
		return err
	}

	err = a.getMeta(path)
	if err != nil {
		return err
	}

	err = a.getEML(path)
	if err != nil {
		return err
	}

	return nil
}

// Meta returns the Meta object of the archive.
func (a *arch) Meta() *meta.Meta {
	return a.metaData
}

// EML returns the EML object of the archive.
func (a *arch) EML() *eml.EML {
	return a.emlData
}

// CoreSlice takes an offset and a limit and returns a slice of slices of
// strings, each slice representing a row of the core file. If limit and
// offset are provided, it returns the corresponding subset of the data.
func (a *arch) CoreSlice(offset, limit int) ([][]string, error) {
	return a.dcFile.CoreData(a.metaData, offset, limit)
}

// CoreStream takes a channel and populates the channel with slices of
// strings, each slice representing a row of the core file. The channel
// is closed when the data is exhausted.
func (a *arch) CoreStream(chCore chan<- []string) error {
	return a.dcFile.CoreStream(a.metaData, chCore)
}

// ExtensionSlice takes an index, offset and limit and returns a slice of
// slices of strings, each slice representing a row of the extension file.
// Index corresponds the index of the extension in the extension list.
// If limit and offset are provided, it returns the corresponding subset
// of the data.
func (a *arch) ExtensionSlice(index, offset, limit int) ([][]string, error) {
	return a.dcFile.ExtensionData(index, a.metaData, offset, limit)
}

// ExtensionStream takes an index and a channel and populates the channel
// with slices of strings, each slice representing a row of the extension
// file. The channel is closed when the data is exhausted.
// Index corresponds the index of the extension in the extension list.
func (a *arch) ExtensionStream(index int, ch chan<- []string) error {
	return a.dcFile.ExtensionStream(index, a.metaData, ch)
}

func (a *arch) getMeta(path string) error {
	metaFile, err := os.Open(filepath.Join(path, "meta.xml"))
	if err != nil {
		return err
	}

	a.metaData, err = meta.New(metaFile)
	if err != nil {
		return err
	}
	return nil
}

func (a *arch) getEML(path string) error {
	emlFileName := "eml.xml"
	if a.metaData.EMLFile != "" {
		emlFileName = a.metaData.EMLFile
	}

	emlFile, err := os.Open(filepath.Join(path, emlFileName))
	if err != nil {
		return err
	}

	a.emlData, err = eml.New(emlFile)
	if err != nil {
		return err
	}

	return nil
}

func (a *arch) Diagnose() (*diagn.Diagnostics, error) {
	cs, exts, err := a.coreSample()
	if err != nil {
		return nil, err
	}
	if cs == nil {
		return nil, errors.New("no data in the core file")
	}

	prs := <-a.gnpPool
	defer func() { a.gnpPool <- prs }()

	return diagn.New(prs, cs, exts), nil
}

func (a *arch) coreSample() (
	[]map[string]string,
	map[string]string,
	error,
) {
	dt, err := a.CoreSlice(0, 1000)
	if err != nil {
		return nil, nil, err
	}
	m := a.metaData.Simplify()
	coreRows := make([]map[string]string, len(dt))
	exts := make(map[string]string)
	for k, v := range m.ExtensionsData {
		exts[k] = strings.ToLower(v.Location)
	}
	for i, row := range dt {
		coreRows[i] = make(map[string]string)
		for j, val := range row {
			if m.CoreData.FieldsIdx[j].Term == "" {
				continue
			}
			coreRows[i][m.CoreData.FieldsIdx[j].Term] = val
		}
	}
	return coreRows, exts, nil
}
