package dcfileio

import (
	"fmt"

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
		fmt.Println("dir absent")
		return gnsys.MakeDir(d.cfg.Path)
	case gnsys.DirEmpty:
		fmt.Println("dir empty")
		return nil
	case gnsys.DirNotEmpty:
		fmt.Println("clean dir")
		return gnsys.CleanDir(d.cfg.Path)
	default:
		fmt.Println("clean dir")
		return &dcfile.ErrDir{}
	}
}
