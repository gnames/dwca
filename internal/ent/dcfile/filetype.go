package dcfile

import (
	"log/slog"
	"strings"
)

type FileType int

const (
	Unknown FileType = iota
	ZIP              // .zip
	TAR              // .tar
	TARGZ            // .tar.gz
	TARXZ            // .tar.xz
	TARBZ2           //.tar.bz2
)

var ftMap = map[FileType]string{
	Unknown: "unknown",
	ZIP:     "zip",
	TAR:     "tar",
	TARGZ:   "tar-gzip",
	TARXZ:   "tar-xz",
	TARBZ2:  "tar-bz2",
}

func (ft FileType) String() string {
	return ftMap[ft]
}

func NewFileType(file string) FileType {
	if file == "" {
		return ZIP
	}
	switch {
	case strings.HasSuffix(file, ".zip"):
		return ZIP
	case strings.HasSuffix(file, ".tar"):
		return TAR
	case strings.HasSuffix(file, ".tar.gz"):
		return TARGZ
	case strings.HasSuffix(file, ".tar.xz"):
		return TARXZ
	case strings.HasSuffix(file, ".tar.bz2"):
		return TARBZ2
	default:
		slog.Info("Unknown file type, trying ZIP", "file", file)
		return ZIP
	}
}
