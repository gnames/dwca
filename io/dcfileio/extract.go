package dcfileio

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/gnames/dwca/ent/dcfile"
	"github.com/gnames/gnsys"
	"github.com/ulikunitz/xz"
)

// extractTar extracts the content of the DwCA tar file to a temporary
// directory.
func (d *dcfileio) extractTar() error {
	// Open the tar archive for reading.
	file, err := os.Open(d.fpath)
	if err != nil {
		return dcfile.ErrExtract{Path: d.fpath, Err: err}
	}
	defer file.Close()

	tr := tar.NewReader(file)
	return d.untar(tr)
}

func (d *dcfileio) untar(tarReader *tar.Reader) error {
	var writer *os.File
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return dcfile.ErrExtract{Path: d.fpath, Err: err}
		}

		// Get the individual filepath from the header.
		filepath := filepath.Join(d.cfg.ExtractPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Handle directory.
			err = os.MkdirAll(filepath, os.FileMode(header.Mode))
			if err != nil {
				return dcfile.ErrExtract{Path: d.fpath, Err: err}
			}
		case tar.TypeReg:
			// Handle regular file.
			writer, err = os.Create(filepath)
			if err != nil {
				return dcfile.ErrExtract{Path: d.fpath, Err: err}
			}
			io.Copy(writer, tarReader)
			writer.Close()
		default:
			return dcfile.ErrExtract{Path: d.fpath, Err: err}
		}
	}
	state := gnsys.GetDirState(d.cfg.ExtractPath)
	if state == gnsys.DirEmpty {
		return dcfile.ErrExtract{
			Path: d.cfg.ExtractPath,
			Err:  errors.New("bad tar file"),
		}
	}
	return nil
}

func (d *dcfileio) extractTarGz() error {
	// Open the .tar.gz archive for reading.
	file, err := os.Open(d.fpath)
	if err != nil {
		return dcfile.ErrExtract{Path: d.fpath, Err: err}
	}
	defer file.Close()

	// Create a new gzip reader.
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return dcfile.ErrExtract{Path: d.fpath, Err: err}
	}
	defer gzReader.Close()

	// Create a new tar reader from the gzip reader.
	tr := tar.NewReader(gzReader)
	return d.untar(tr)
}

func (d *dcfileio) extractTarBz2() error {
	// Open the .tar.gz archive for reading.
	file, err := os.Open(d.fpath)
	if err != nil {
		return dcfile.ErrExtract{Path: d.fpath, Err: err}
	}
	defer file.Close()

	// Create a new bz2 reader.
	bzReader := bzip2.NewReader(file)

	// Create a new tar reader from the gzip reader.
	tr := tar.NewReader(bzReader)
	return d.untar(tr)
}

func (d *dcfileio) extractTarXz() error {
	// Open the .tar.gz archive for reading.
	file, err := os.Open(d.fpath)
	if err != nil {
		return dcfile.ErrExtract{Path: d.fpath, Err: err}
	}
	defer file.Close()

	xzReader, err := xz.NewReader(file)
	if err != nil {
		return dcfile.ErrExtract{Path: d.fpath, Err: err}
	}

	// Create a new tar reader from the gzip reader.
	tr := tar.NewReader(xzReader)
	return d.untar(tr)
}

func (d *dcfileio) extractZip() error {
	// Open the zip file for reading.
	r, err := zip.OpenReader(d.fpath)
	if err != nil {
		return dcfile.ErrExtract{Path: d.fpath, Err: err}
	}
	defer r.Close()

	for _, f := range r.File {
		// Construct the full path for the file/directory and ensure its directory exists.
		fpath := filepath.Join(d.cfg.ExtractPath, f.Name)
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return dcfile.ErrExtract{Path: fpath, Err: err}
		}

		// If it's a directory, move on to the next entry.
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Open the file within the zip.
		rc, err := f.Open()
		if err != nil {
			return dcfile.ErrExtract{Path: fpath, Err: err}
		}
		defer rc.Close()

		// Create a file in the filesystem.
		outFile, err := os.OpenFile(
			fpath,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			f.Mode(),
		)
		if err != nil {
			return dcfile.ErrExtract{Path: fpath, Err: err}
		}
		defer outFile.Close()

		// Copy the contents of the file from the zip to the new file.
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return dcfile.ErrExtract{Path: fpath, Err: err}
		}
	}

	return nil
}
