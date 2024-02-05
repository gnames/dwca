package dcfileio

import (
	"os"
	"path/filepath"

	"github.com/gnames/dwca/config"
	"github.com/gnames/dwca/internal/ent/dcfile"
	"github.com/gnames/gnsys"
)

type dcfileio struct {
	// config
	cfg config.Config
	// file type
	fileType dcfile.FileType
	// filePath is the DwC Archive file path
	filePath string
	// arcPath is the path where all DwCA data files are located.
	arcPath string
}

// New creates a new DCFile object.
func New(cfg config.Config, path string) (dcfile.DCFile, error) {
	exists, _ := gnsys.FileExists(path)
	if !exists {
		return nil, &dcfile.ErrFileNotFound{Path: path}
	}
	res := &dcfileio{
		cfg:      cfg,
		filePath: path,
		fileType: dcfile.NewFileType(path),
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
	switch d.fileType {
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
		return &dcfile.ErrUnknownArchiveType{FileType: d.fileType}
	}
}

func (d *dcfileio) ArchiveDir() (string, error) {
	if d.arcPath != "" {
		return d.arcPath, nil
	}
	var dirs []string
	err := filepath.Walk(d.cfg.ExtractPath,
		func(path string, info os.FileInfo, err error,
		) error {
			if err != nil {
				return err // handle the error and possibly abort the Walk
			}

			// Check if the current path is the file we're looking for
			if !info.IsDir() && info.Name() == "meta.xml" {
				dir := filepath.Dir(path) // get the directory of the file
				dirs = append(dirs, dir)  // add it to the slice
			}

			return nil
		})

	if err != nil {
		return "", err
	}

	if len(dirs) == 0 {
		return "", &dcfile.ErrMetaFileNotFound{}
	}

	if len(dirs) > 1 {
		return "", &dcfile.ErrMultipleMetaFiles{}
	}

	return dirs[0], nil
}

func (d *dcfileio) Close() error {
	err := os.RemoveAll(d.cfg.ExtractPath)
	if err != nil {
		return err
	}
	return os.RemoveAll(d.cfg.DownloadPath)
}
