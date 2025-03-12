package impio

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gnames/dwca"
	"github.com/gnames/gnsys"
)

func (_ *impio) Import(srcPath, dstDir string) (string, error) {
	var err error
	var isURL bool
	var rootDir string

	slog.Info("Importing data from DwCA file")
	srcPath, isURL, err = download(srcPath)
	if err != nil {
		return "", err
	}
	if isURL {
		defer os.RemoveAll(filepath.Dir(srcPath))
	}

	err = extract(srcPath, dstDir)
	if err != nil {
		return "", err
	}

	rootDir, err = getRootDir(dstDir)
	if err != nil {
		return "", err
	}

	return rootDir, nil
}

func download(srcPath string) (string, bool, error) {
	var err error
	var dlDir string
	if !strings.HasPrefix(srcPath, "http") {
		return srcPath, false, nil
	}

	slog.Info("Downloading DwCA file", "url", srcPath)
	dlDir, err = os.MkdirTemp("", "dwca-download")
	if err != nil {
		err = &dwca.ErrDownload{URL: srcPath, Err: err}
		return "", true, err
	}

	// srcPath is now local
	srcPath, err = gnsys.Download(srcPath, dlDir, true)
	if err != nil {
		err = &dwca.ErrDownload{URL: srcPath, Err: err}
		return "", true, err
	}

	return srcPath, true, nil
}

func extract(srcPath, dstDir string) error {
	var err error
	ft := gnsys.GetFileType(srcPath)
	switch ft {
	case gnsys.ZipFT:
		err = gnsys.ExtractZip(srcPath, dstDir)
	case gnsys.TarGzFT:
		err = gnsys.ExtractTarGz(srcPath, dstDir)
	default:
		err = errors.New("unknown file type")
		return &dwca.ErrExtractArchive{File: srcPath, Err: err}
	}
	if err != nil {
		return &dwca.ErrExtractArchive{File: srcPath, Err: err}
	}

	return nil
}

// getRootDir determines the directory where the files of DarwinCore archive
// reside.
func getRootDir(path string) (string, error) {
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
		return "", &dwca.ErrNoMetaFile{}
	}

	if len(dirs) > 1 {
		return "", &dwca.ErrMultipleMetaFiles{}
	}

	return dirs[0], nil
}
