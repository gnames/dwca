package dwca

import (
	"context"
	"log/slog"
	"path/filepath"

	"github.com/gnames/dwca/pkg/ent/meta"
	"golang.org/x/sync/errgroup"
)

func (a *arch) processExtensionsOutput() error {
	for i := range a.meta.Extensions {
		a.processExt(i)
	}
	return nil
}

func (a *arch) processExt(idx int) {
	ext := a.meta.Extensions[idx]
	extType := ext.RowType
	extType = filepath.Base(extType)
	slog.Info("Processing extension", "ext", extType)
	var maxIdx int
	for _, v := range ext.Fields {
		if v.Idx > maxIdx {
			maxIdx = v.Idx
		}
	}
	a.updateOutputMetaExt(idx)

	chIn := make(chan []string)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return a.saveExtOutput(ctx, idx, chIn)
	})

	_, err := a.ExtensionStream(ctx, idx, chIn)
	if err != nil {
		slog.Error(
			"Error processing extension",
			"file", ext.Files.Location,
			"error", err,
		)
	}
}

func (a *arch) saveExtOutput(
	ctx context.Context,
	idx int,
	chIn <-chan []string,
) error {
	ext := a.outputMeta.Extensions[idx]
	file := ext.Files.Location
	fields := meta.Headers(idx, ext.Fields)
	delim := ext.FieldsTerminatedBy
	return a.dcFile.ExportCSVStream(ctx, file, fields, delim, chIn)
}

func (a *arch) updateOutputMetaExt(idx int) {
	ext := a.outputMeta.Extensions[idx]

	file := ext.Files.Location
	e := filepath.Ext(file)
	location := file[:len(file)-len(e)] + ".txt"

	delim := ","
	if a.cfg.OutputCSVType == "tsv" {
		delim = `\t`
	}

	ext.Files.Location = location
	ext.FieldsTerminatedBy = delim
	ext.LinesTerminatedBy = `\n`
	ext.FieldsEnclosedBy = `"`
	ext.IgnoreHeaderLines = "1"
}
