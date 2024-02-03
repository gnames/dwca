package dcfileio

import (
	"github.com/gnames/dwca/config"
	"github.com/gnames/dwca/ent/dcfile"
	"github.com/gnames/gnsys"
)

type dcfileio struct {
	// config
	cfg config.Config
	// file type
	ft dcfile.FileType
	// fpath is the DwC Archive file path
	fpath string
}

// New creates a new DCFile object.
func New(cfg config.Config, path string) (dcfile.DCFile, error) {
	exists, _ := gnsys.FileExists(path)
	if !exists {
		return nil, dcfile.ErrFileNotFound{Path: path}
	}
	res := &dcfileio{
		cfg:   cfg,
		fpath: path,
		ft:    dcfile.NewFileType(path),
	}
	return res, nil
}

func (d *dcfileio) Init() error {
	err := d.touchDirs()
	if err != nil {
		return err
	}
	return nil
}

func (d *dcfileio) Extract() error {
	switch d.ft {
	case dcfile.TAR:
		return d.extractTar()
	case dcfile.TARGZ:
		return d.extractTarGz()
	case dcfile.TARBZ2:
		return d.extractTarBz2()
	case dcfile.TARXZ:
		return d.extractTarXz()
	case dcfile.ZIP:
		return d.extractZip()
	default:
		return dcfile.ErrUnknownFileType{FileType: d.ft}
	}
}

func (d *dcfileio) Close() error {
	return nil
}
