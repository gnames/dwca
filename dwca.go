package dwca

import (
	"github.com/gnames/dwca/config"
	"github.com/gnames/dwca/ent/meta"
	"github.com/gnames/dwca/internal/ent/dcfile"
)

type arch struct {
	cfg config.Config
	df  dcfile.DCFile
	md  *meta.Meta
}

func New(cfg config.Config, df dcfile.DCFile) Archive {
	return &arch{cfg: cfg, df: df}
}

func (a *arch) Load() error {
	err := a.df.Extract()
	if err != nil {
		return err
	}
	path, err := a.df.ArchiveDir()
	if err != nil {
		return err
	}
	_ = path
	return nil
}

func (a *arch) Meta() *meta.Meta {
	return a.md
}

func (a *arch) Config() config.Config {
	return a.cfg
}
