package dcfileio

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gnames/dwca/internal/ent"
	"github.com/gnames/dwca/internal/ent/dcfile"
	"github.com/gnames/dwca/internal/io/factory"
	"github.com/gnames/dwca/pkg/config"
	"github.com/gnames/dwca/pkg/ent/meta"
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
	if !exists && path != "" && !strings.HasPrefix(path, "http") {
		return nil, &dcfile.ErrFileNotFound{Path: path}
	}
	res := &dcfileio{
		cfg:      cfg,
		filePath: path,
		fileType: dcfile.NewFileType(path),
	}
	return res, nil
}

// ResetTempDirs creates empty filesystem structure for the DwCA archive.
func (d *dcfileio) ResetTempDirs() error {
	err := d.resetDirs()
	if err != nil {
		return err
	}
	return nil
}

func (d *dcfileio) SetFilePath(path string) {
	d.filePath = path
	d.fileType = dcfile.NewFileType(path)
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

// ArchiveDir determines the directory where the files of DarwinCore archive
// reside.
func (d *dcfileio) ArchiveDir(path string) (string, error) {
	var dirs []string
	err := filepath.Walk(path,
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

func (d *dcfileio) CoreData(
	root string,
	meta *meta.Meta,
	offset, limit int,
) ([][]string, error) {
	if meta == nil {
		return nil, &dcfile.ErrCoreRead{Err: errors.New("*meta.Meta is nil")}
	}

	path, err := d.basePath(root, meta.Core.Files.Location)
	if err != nil {
		return nil, err
	}

	attr := ent.CSVAttr{
		Path:             path,
		ColSep:           colSep(meta.Core.FieldsTerminatedBy),
		Quote:            meta.Core.FieldsEnclosedBy,
		IgnoreHeader:     meta.Core.IgnoreHeaderLines,
		BadRowProcessing: d.cfg.WrongFieldsNum,
	}

	r, err := factory.CSVReader(attr)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	res, err := r.ReadSlice(offset, limit)
	if err != nil {
		return nil, &dcfile.ErrCoreRead{Err: err}
	}

	return res, nil
}

func (d *dcfileio) CoreStream(
	ctx context.Context,
	root string,
	meta *meta.Meta,
	coreChan chan<- []string,
) (int, error) {
	path, err := d.basePath(root, meta.Core.Files.Location)
	if err != nil {
		return 0, err
	}
	attr := ent.CSVAttr{
		Path:             path,
		ColSep:           colSep(meta.Core.FieldsTerminatedBy),
		Quote:            meta.Core.FieldsEnclosedBy,
		IgnoreHeader:     meta.Core.IgnoreHeaderLines,
		BadRowProcessing: d.cfg.WrongFieldsNum,
	}

	r, err := factory.CSVReader(attr)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	count, err := r.Read(ctx, coreChan)
	close(coreChan)

	if err != nil {
		return int(count), &dcfile.ErrExtensionRead{Err: err}
	}
	slog.Info("Processed core", "lines", humanize.Comma(int64(count)))

	return int(count), nil
}

func (d *dcfileio) ExtensionData(
	index int, root string, meta *meta.Meta, offset, limit int,
) ([][]string, error) {
	if meta == nil {
		return nil, &dcfile.ErrExtensionRead{Err: errors.New("*meta.Meta is nil")}
	}
	if len(meta.Extensions) <= index {
		return nil, &dcfile.ErrExtensionRead{Err: errors.New("index out of range")}
	}

	ext := meta.Extensions[index]
	path, err := d.basePath(root, ext.Files.Location)
	if err != nil {
		return nil, err
	}

	attr := ent.CSVAttr{
		Path:             path,
		ColSep:           colSep(ext.FieldsTerminatedBy),
		IgnoreHeader:     ext.IgnoreHeaderLines,
		BadRowProcessing: d.cfg.WrongFieldsNum,
	}

	r, err := factory.CSVReader(attr)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	res, err := r.ReadSlice(offset, limit)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *dcfileio) ExtensionStream(
	ctx context.Context,
	index int,
	root string,
	meta *meta.Meta,
	extChan chan<- []string,
) (int, error) {
	if meta == nil {
		return 0, &dcfile.ErrExtensionRead{Err: errors.New("*meta.Meta is nil")}
	}
	if len(meta.Extensions) <= index {
		return 0, &dcfile.ErrExtensionRead{Err: errors.New("index out of range")}
	}
	ext := meta.Extensions[index]
	extType := ext.RowType
	extType = filepath.Base(extType)

	path, err := d.basePath(root, ext.Files.Location)
	if err != nil {
		return 0, err
	}

	attr := ent.CSVAttr{
		Path:             path,
		ColSep:           colSep(ext.FieldsTerminatedBy),
		IgnoreHeader:     ext.IgnoreHeaderLines,
		BadRowProcessing: d.cfg.WrongFieldsNum,
	}

	r, err := factory.CSVReader(attr)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	count, err := r.Read(ctx, extChan)
	close(extChan)

	if err != nil {
		return int(count), &dcfile.ErrExtensionRead{Err: err}
	}

	slog.Info(
		"Processed %s",
		"lines", humanize.Comma(int64(count)), "ext", extType,
	)
	return int(count), nil
}

func (d *dcfileio) ExportCSVStream(
	ctx context.Context,
	file string,
	headers []string,
	delim string,
	outChan <-chan []string,
) error {
	attr := ent.CSVAttr{
		Headers:          headers,
		Path:             filepath.Join(d.cfg.OutputPath, file),
		ColSep:           colSep(delim),
		BadRowProcessing: d.cfg.WrongFieldsNum,
	}
	w, err := factory.CSVWriter(attr)
	if err != nil {
		return err
	}
	defer w.Close()

	err = w.Write(ctx, outChan)
	if err != nil {
		return err
	}

	return nil
}

// 	f, err := os.Create(path)
// 	if err != nil {
// 		return &dcfile.ErrSaveCSV{Err: err}
// 	}
// 	defer f.Close()
//
// 	w := csv.NewWriter(f)
// 	w.Comma = ','
// 	if delim == "\\t" {
// 		w.Comma = '\t'
// 	}
//
// 	err = w.Write(headers)
// 	if err != nil {
// 		return err
// 	}
// 	for row := range outChan {
// 		err = w.Write(row)
// 		if err != nil {
// 			for range outChan {
// 			}
// 			return err
// 		}
// 		select {
// 		case <-ctx.Done():
// 			return &dcfile.ErrContext{Err: ctx.Err()}
// 		default:
// 		}
// 	}
// 	w.Flush()
// 	return nil
// }

func (d *dcfileio) SaveToFile(fileName string, bs []byte) error {
	path := filepath.Join(d.cfg.OutputPath, fileName)
	return os.WriteFile(path, bs, 0644)
}

func (d *dcfileio) Zip(inputDir, fileZip string) error {
	w, err := os.Create(fileZip)
	if err != nil {
		return err
	}
	defer w.Close()

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	filepath.WalkDir(inputDir,
		func(path string, e os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if e.IsDir() {
				return nil // Skip directories
			}

			relPath, err := filepath.Rel(d.cfg.OutputPath, path)
			if err != nil {
				return err
			}

			fileInfo, err := e.Info()
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(fileInfo)
			if err != nil {
				return err
			}

			header.Name = relPath // Store relative path within the zip
			header.Method = zip.Deflate

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		})

	return nil
}

func (d *dcfileio) TarGz(inputDir, fileTar string) error {
	w, err := os.Create(fileTar)
	if err != nil {
		return err
	}
	defer w.Close()

	gzWriter := gzip.NewWriter(w)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	filepath.WalkDir(inputDir,
		func(path string, de os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if de.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(d.cfg.OutputPath, path)
			if err != nil {
				return err
			}

			fileInfo, err := de.Info()
			if err != nil {
				return err
			}

			header, err := tar.FileInfoHeader(fileInfo, relPath)
			if err != nil {
				return err
			}

			err = tarWriter.WriteHeader(header)
			if err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarWriter, file)
			return err
		})
	return nil
}

func (d *dcfileio) Close() error {
	err := os.RemoveAll(d.cfg.ExtractPath)
	if err != nil {
		return err
	}
	err = os.RemoveAll(d.cfg.DownloadPath)
	if err != nil {
		return err
	}
	return os.RemoveAll(d.cfg.OutputPath)
}

func colSep(s string) rune {
	if s == "\\t" {
		return '\t'
	}

	if len(s) > 0 {
		return rune(s[0])
	}
	return ','
}

func (d *dcfileio) basePath(root, filePath string) (string, error) {
	path, err := d.ArchiveDir(root)
	if err != nil {
		return "", err
	}
	path = filepath.Join(path, filePath)
	return path, nil
}
