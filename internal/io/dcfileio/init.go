package dcfileio

import (
	"github.com/gnames/dwca/internal/ent/dcfile"
	"github.com/gnames/gnsys"
)

func (d *dcfileio) touchDirs() error {
	err := d.rootDir()
	if err != nil {
		return err
	}

	err = gnsys.MakeDir(d.cfg.DownloadPath)
	if err != nil {
		return err
	}

	err = gnsys.MakeDir(d.cfg.ExtractPath)
	if err != nil {
		return err
	}
	return nil
}

func (d *dcfileio) rootDir() error {
	switch gnsys.GetDirState(d.cfg.Path) {
	case gnsys.DirAbsent:
		return gnsys.MakeDir(d.cfg.Path)
	case gnsys.DirEmpty:
		return nil
	case gnsys.DirNotEmpty:
		return gnsys.CleanDir(d.cfg.Path)
	default:
		return &dcfile.ErrDir{}
	}
}
