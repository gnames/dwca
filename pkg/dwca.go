package dwca

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gnames/dwca/internal/ent/dcfile"
	"github.com/gnames/dwca/internal/ent/diagn"
	"github.com/gnames/dwca/pkg/config"
	"github.com/gnames/dwca/pkg/ent/eml"
	"github.com/gnames/dwca/pkg/ent/meta"
	"github.com/gnames/gnlib/ent/gnvers"
	"github.com/gnames/gnparser"
)

// arch implements Archive interface.
type arch struct {
	// root contains the path where DWCA uncompressed files and directories
	// are located. For example it can be cfg.ExtractPath or cfg.OutputPath.
	root string

	// cfg is the configuration object of the archive.
	cfg config.Config

	// dcFile is the object that handles the DwCA archive filesystem operations.
	dcFile dcfile.DCFile

	// meta is the object that holds the metadata of the DwCA archive.
	meta *meta.Meta

	// outputMeta is the object that holds the metadata of the DwCA archive and
	// is modified for output DwCA file.
	outputMeta *meta.Meta

	// metaSimple is the simplified version of the metadata. It is useful to
	// access metadata fields by their names or indices.
	metaSimple *meta.MetaSimple

	// emlData is the object that holds the EML data of the DwCA archive.
	emlData *eml.EML

	// dgn is the object that holds the diagnostics that detect semantically
	// fuzzy fields and how are they used in the DwCA file.
	dgn *diagn.Diagnostics

	// gnpPool is a pool of GNparser objects to be used in parallel processing.
	gnpPool chan gnparser.GNparser

	// taxon contains information about DarwinCore fields that are relevant
	// for taxon information.
	taxon *taxon

	// hierarchy is used when core contains parent-child relationship to
	// represent a hierarchy.
	hierarchy map[string]*hNode
}

// New creates a new Archive object. It takes configuration file and necessary
// internal objects to handle the DwCA archive.
func New(cfg config.Config, df dcfile.DCFile) Archive {
	res := &arch{cfg: cfg, dcFile: df}
	poolSize := cfg.JobsNum
	gnpPool := make(chan gnparser.GNparser, poolSize)
	for i := 0; i < poolSize; i++ {
		cfgGNP := gnparser.NewConfig()
		gnpPool <- gnparser.New(cfgGNP)
	}
	res.gnpPool = gnpPool
	res.hierarchy = make(map[string]*hNode)
	return res
}

func Version() gnvers.Version {
	return gnvers.Version{
		Version: Vers,
		Build:   Build,
	}
}

// Config returns the configuration object of the archive.
func (a *arch) Config() config.Config {
	return a.cfg
}

// Load extracts the archive and loads data for EML and Meta.
func (a *arch) Load(path string) error {
	var err error
	slog.Info("Loading data from input DwCA file")

	a.root = path

	if a.root == a.cfg.ExtractPath {
		err = a.dcFile.Extract()
		if err != nil {
			return err
		}
	}
	path, err = a.dcFile.ArchiveDir(path)
	if err != nil {
		return err
	}

	slog.Info("Reading meta.xml and eml.xml files")
	err = a.getMeta(path)
	if err != nil {
		return err
	}

	a.metaSimple = a.meta.Simplify()

	err = a.getEML(path)
	if err != nil {
		return err
	}

	slog.Info("Analyzing the archive")
	err = a.getDiagnostics()
	if err != nil {
		return err
	}

	return nil
}

// Close closes the archive and all associated files.
func (a *arch) Close() error {
	return a.dcFile.Close()
}

// Meta returns the Meta object of the archive.
func (a *arch) Meta() *meta.Meta {
	return a.meta
}

// EML returns the EML object of the archive.
func (a *arch) EML() *eml.EML {
	return a.emlData
}

// CoreSlice takes an offset and a limit and returns a slice of slices of
// strings, each slice representing a row of the core file. If limit and
// offset are provided, it returns the corresponding subset of the data.
func (a *arch) CoreSlice(offset, limit int) ([][]string, error) {
	return a.dcFile.CoreData(a.root, a.meta, offset, limit)
}

// CoreStream takes a channel and populates the channel with slices of
// strings, each slice representing a row of the core file. The channel
// is closed when the data is exhausted.
func (a *arch) CoreStream(
	ctx context.Context,
	chCore chan<- []string,
) (int, error) {
	return a.dcFile.CoreStream(ctx, a.root, a.meta, chCore)
}

// ExtensionSlice takes an index, offset and limit and returns a slice of
// slices of strings, each slice representing a row of the extension file.
// Index corresponds the index of the extension in the extension list.
// If limit and offset are provided, it returns the corresponding subset
// of the data.
func (a *arch) ExtensionSlice(index, offset, limit int) ([][]string, error) {
	return a.dcFile.ExtensionData(index, a.root, a.meta, offset, limit)
}

// ExtensionStream takes an index and a channel and populates the channel
// with slices of strings, each slice representing a row of the extension
// file. The channel is closed when the data is exhausted.
// Index corresponds the index of the extension in the extension list.
func (a *arch) ExtensionStream(
	ctx context.Context,
	index int,
	ch chan<- []string,
) (int, error) {
	return a.dcFile.ExtensionStream(ctx, index, a.root, a.meta, ch)
}

func (a *arch) getMeta(path string) error {
	metaFile, err := os.Open(filepath.Join(path, "meta.xml"))
	if err != nil {
		return err
	}
	defer metaFile.Close()

	a.meta, err = meta.New(metaFile)
	if err != nil {
		return err
	}

	// rewind file back to the beginning
	_, err = metaFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	a.outputMeta, err = meta.New(metaFile)
	if err != nil {
		return err
	}

	return nil
}

func (a *arch) getEML(path string) error {
	emlFileName := "eml.xml"
	if a.meta.EMLFile != "" {
		emlFileName = a.meta.EMLFile
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

func (a *arch) getDiagnostics() error {
	cs, exts, err := a.coreSample()
	if err != nil {
		return err
	}
	if cs == nil {
		return errors.New("no data in the core file")
	}

	prs := <-a.gnpPool
	defer func() { a.gnpPool <- prs }()

	res := diagn.New(prs, cs, exts)
	a.dgn = res
	return nil
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
	m := a.meta.Simplify()
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

func (a *arch) Normalize() error {
	slog.Info("Processing Core")
	err := a.processCoreOutput()
	if err != nil {
		return err
	}

	slog.Info("Processing Extensions")
	err = a.processExtensionsOutput()
	if err != nil {
		return err
	}

	slog.Info("Saving normalized meta.xml and eml.xml files")
	err = a.saveMetaOutput()
	if err != nil {
		return err
	}

	err = a.saveEmlOutput()
	if err != nil {
		return err
	}

	return nil
}

func (a *arch) ZipNormalized(filePath string) error {
	slog.Info("Creating zip archive", "output", filePath)
	err := a.dcFile.Zip(a.cfg.OutputPath, filePath)
	if err != nil {
		return err
	}
	slog.Info("The zip archive created", "output", filePath)
	return nil
}

func (a *arch) TarGzNormalized(filePath string) error {
	slog.Info("Creating tar.gz archive", "output", filePath)
	err := a.dcFile.TarGz(a.cfg.OutputPath, filePath)
	if err != nil {
		return err
	}
	slog.Info("The tar.gz archive created", "output", filePath)
	return nil
}
