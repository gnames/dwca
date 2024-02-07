package dwca

import (
	"os"
	"path/filepath"

	"github.com/gnames/dwca/config"
	"github.com/gnames/dwca/ent/eml"
	"github.com/gnames/dwca/ent/meta"
	"github.com/gnames/dwca/internal/ent/dcfile"
)

type arch struct {
	cfg      config.Config
	dcFile   dcfile.DCFile
	metaData *meta.Meta
	emlData  *eml.EML
}

func New(cfg config.Config, df dcfile.DCFile) Archive {
	return &arch{cfg: cfg, dcFile: df}
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

// CoreData takes an offset and a limit and returns a slice of slices of
// strings, each slice representing a row of the core file. If limit and
// offset are provided, it returns the corresponding subset of the data.
func (a *arch) CoreData(offset, limit int) ([][]string, error) {
	return a.dcFile.CoreData(a.metaData, offset, limit)
}

// CoreStream takes a channel and populates the channel with slices of
// strings, each slice representing a row of the core file. The channel
// is closed when the data is exhausted.
func (a *arch) CoreStream(chCore chan<- []string) error {
	return a.dcFile.CoreStream(a.metaData, chCore)
}

// ExtensionData takes an index, offset and limit and returns a slice of
// slices of strings, each slice representing a row of the extension file.
// Index corresponds the index of the extension in the extension list.
// If limit and offset are provided, it returns the corresponding subset
// of the data.
func (a *arch) ExtensionData(index, offset, limit int) ([][]string, error) {
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
