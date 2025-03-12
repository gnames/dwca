package dwcio

import (
	"github.com/gnames/dwca"
	"github.com/gnames/dwca/config"
)

type dwcio struct {
	cfg        config.Config
	archiveDir string
	*dwca.Meta
}

func New(cfg config.Config) dwca.Archive {
	res := dwcio{
		cfg: cfg,
	}
	return &res
}

func (d *dwcio) Config() config.Config {
	return d.cfg
}
