package dcfileio

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/gnames/dwca/internal/ent/dcfile"
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
	err := d.resetDirs()
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

func (d *dcfileio) CoreData(
	meta *meta.Meta,
	offset, limit int,
) ([][]string, error) {
	if meta == nil {
		return nil, &dcfile.ErrCoreRead{Err: errors.New("*meta.Meta is nil")}
	}

	attr := fileAttrs{
		path:         meta.Core.Files.Location,
		colSep:       meta.Core.FieldsTerminatedBy,
		ignoreHeader: meta.Core.IgnoreHeaderLines,
	}

	r, f, err := d.openCSV(attr)
	if err != nil {
		return nil, &dcfile.ErrCoreRead{Err: err}
	}
	defer f.Close()

	// ignore headers if they are given
	if attr.ignoreHeader == "1" {
		r.Read()
	}
	var res [][]string

	var count int
	for {
		count++

		if limit > 0 && len(res) == limit {
			break
		}

		row, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, &dcfile.ErrCoreRead{Err: err}
		}

		if offset > 0 && count <= offset {
			continue
		}
		res = append(res, row)
	}

	return res, nil
}

func (d *dcfileio) CoreStream(ctx context.Context, meta *meta.Meta, coreChan chan<- []string) error {
	attr := fileAttrs{
		path:         meta.Core.Files.Location,
		colSep:       meta.Core.FieldsTerminatedBy,
		ignoreHeader: meta.Core.IgnoreHeaderLines,
	}

	r, f, err := d.openCSV(attr)
	if err != nil {
		return &dcfile.ErrCoreRead{Err: err}
	}
	defer f.Close()

	// ignore headers if they are given
	if attr.ignoreHeader == "1" {
		r.Read()
	}

loop:
	for {
		row, err := r.Read()
		if err == io.EOF {
			break loop
		}
		if err != nil {
			return &dcfile.ErrCoreRead{Err: err}
		}

		select {
		case <-ctx.Done():
			return &dcfile.ErrContext{Err: ctx.Err()}
		default:
			coreChan <- row
		}
	}

	close(coreChan)
	return nil
}

func (d *dcfileio) ExtensionData(
	index int, meta *meta.Meta, offset, limit int,
) ([][]string, error) {
	if meta == nil {
		return nil, &dcfile.ErrExtensionRead{Err: errors.New("*meta.Meta is nil")}
	}
	if len(meta.Extensions) <= index {
		return nil, &dcfile.ErrExtensionRead{Err: errors.New("index out of range")}
	}
	ext := meta.Extensions[index]

	attr := fileAttrs{
		path:         ext.Files.Location,
		colSep:       ext.FieldsTerminatedBy,
		ignoreHeader: ext.IgnoreHeaderLines,
	}

	r, f, err := d.openCSV(attr)
	if err != nil {
		return nil, &dcfile.ErrExtensionRead{Err: err}
	}
	defer f.Close()

	// ignore headers if they are given
	if attr.ignoreHeader == "1" {
		r.Read()
	}
	var res [][]string

	var count int
	for {
		count++

		if limit > 0 && len(res) == limit {
			break
		}

		row, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, &dcfile.ErrCoreRead{Err: err}
		}

		if offset > 0 && count <= offset {
			continue
		}
		res = append(res, row)
	}

	return res, nil
}

func (d *dcfileio) ExtensionStream(
	ctx context.Context,
	index int,
	meta *meta.Meta,
	extChan chan<- []string,
) error {
	if meta == nil {
		return &dcfile.ErrExtensionRead{Err: errors.New("*meta.Meta is nil")}
	}
	if len(meta.Extensions) <= index {
		return &dcfile.ErrExtensionRead{Err: errors.New("index out of range")}
	}
	ext := meta.Extensions[index]

	attr := fileAttrs{
		path:         ext.Files.Location,
		colSep:       ext.FieldsTerminatedBy,
		ignoreHeader: ext.IgnoreHeaderLines,
	}

	r, f, err := d.openCSV(attr)
	if err != nil {
		return &dcfile.ErrExtensionRead{Err: err}
	}
	defer f.Close()
	// ignore headers if they are given
	if attr.ignoreHeader == "1" {
		r.Read()
	}

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return &dcfile.ErrExtensionRead{Err: err}
		}

		select {
		case <-ctx.Done():
			return &dcfile.ErrContext{Err: ctx.Err()}
		default:
			extChan <- row
		}
	}

	close(extChan)
	return nil
}

func (d *dcfileio) ExportCSVStream(
	ctx context.Context,
	file string,
	fields []string,
	outChan <-chan []string,
) error {
	path := filepath.Join(d.cfg.OutputPath, file)
	f, err := os.Create(path)
	if err != nil {
		return &dcfile.ErrSaveCSV{Err: err}
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = ','

	err = w.Write(fields)
	if err != nil {
		return err
	}

	for row := range outChan {
		select {
		case <-ctx.Done():
			return &dcfile.ErrContext{Err: ctx.Err()}
		default:
			err = w.Write(row)
			if err != nil {
				return err
			}
		}
	}
	w.Flush()
	return nil
}

func (d *dcfileio) SaveToFile(fileName string, bs []byte) error {
	path := filepath.Join(d.cfg.OutputPath, fileName)
	return os.WriteFile(path, bs, 0644)
}

func (d *dcfileio) ZipOutput(filePath string) error {
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	filepath.WalkDir(d.cfg.OutputPath,
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

func (dd *dcfileio) TarGzOutput(tarGzFilename string) error {
	w, err := os.Create(tarGzFilename)
	if err != nil {
		return err
	}
	defer w.Close()

	gzWriter := gzip.NewWriter(w)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	filepath.WalkDir(dd.cfg.OutputPath,
		func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(dd.cfg.OutputPath, path)
			if err != nil {
				return err
			}

			fileInfo, err := d.Info()
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

type fileAttrs struct {
	path         string
	colSep       string
	ignoreHeader string
}

func (d *dcfileio) openCSV(attr fileAttrs) (*csv.Reader, *os.File, error) {
	path := attr.path
	if path == "" {
		return nil, nil, errors.New("core file location is empty")
	}
	basePath, err := d.ArchiveDir()
	if err != nil {
		return nil, nil, err
	}

	colSep := ','
	if attr.colSep == "\\t" {
		colSep = '\t'
	}

	path = filepath.Join(basePath, path)

	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	r := csv.NewReader(f)
	r.Comma = colSep
	// allow variable number of fields
	r.FieldsPerRecord = -1

	if r.Comma == '\t' {
		// lax quotes for tab-separated files
		r.LazyQuotes = true
	}

	return r, f, nil
}
