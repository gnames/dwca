package dwca

import (
	"log/slog"

	"github.com/gnames/dwca/internal/io/dcfileio"
	"github.com/gnames/dwca/pkg/config"
)

// Factory creates a new DwCA object. It takes a list of options for the
// configuration, and a path to the DwCA file. The path is used to initialize
// the DwCA object, and the options are used to configure the object.
// This function is the only place where concrete IO objects are allowed.
func Factory(fpath string, opts ...config.Option) (Archive, error) {
	slog.Info("Creating empty DwCA object", "input", fpath)
	cfg := config.New(opts...)
	dcf, err := dcfileio.New(cfg, fpath)
	if err != nil {
		return nil, err
	}
	err = dcf.Init()
	if err != nil {
		return nil, err
	}

	res := New(cfg, dcf)
	return res, nil
}
