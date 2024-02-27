package dwca

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

func (a *arch) processExtensionsOutput() error {
	for i := range a.meta.Extensions {
		err := a.processExt(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *arch) processExt(idx int) error {
	ext := a.meta.Extensions[idx]
	slog.Info("Processing extension", "file", ext.Files.Location)
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

	err := a.ExtensionStream(ctx, idx, chIn)
	if err != nil {
		return err
	}
	return nil
}

func (a *arch) saveExtOutput(
	ctx context.Context,
	idx int,
	chIn <-chan []string,
) error {
	ext := a.meta.Extensions[idx]
	file := ext.Files.Location
	fields := fieldNames(ext.Fields)
	return a.dcFile.ExportCSVStream(ctx, file, fields, chIn)
}

func (a *arch) updateOutputMetaExt(idx int) {
	ext := a.outputMeta.Extensions[idx]
	ext.FieldsTerminatedBy = ","
	ext.LinesTerminatedBy = "\n"
	ext.FieldsEnclosedBy = "\""
	ext.IgnoreHeaderLines = "1"
}
