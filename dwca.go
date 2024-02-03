package dwca

import (
	"github.com/gnames/dwca/config"
	"github.com/gnames/dwca/ent/dcfile"
)

type arch struct {
	cfg config.Config
	df  dcfile.DCFile
}

func New(cfg config.Config, df dcfile.DCFile) Archive {
	return &arch{cfg: cfg, df: df}
}

func (d *arch) Extract() error {
	return d.df.Extract()
}
