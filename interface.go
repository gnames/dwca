package dwca

import "github.com/gnames/dwca/config"

// Archive defines methods to manage data from DwCA file. It makes it possible
// to either import existing DwCA or create new DwCA file.
type Archive interface {
	Importer
	Accessor
}

// Importer defines methods for importing and ingesting Darwin Core Archive
// (DwCA) files. It provides functionality to download, extract, and parse DwCA
// files, ultimately creating Meta and EML objects that represent the archive's
// structure and metadata.
//
// The Importer interface is responsible for:
//   - Downloading DwCA files from URLs or using local paths.
//   - Extracting the contents of compressed DwCA files (zip, tar.gz).
//   - Parsing the `meta.xml` file to create a `Meta` object.
//   - Parsing the `eml.xml` file to create an `EML` object.
//
// `Meta` objects contain information about the structure and meaning of the
// data files within the archive, including details about the fields and their
// definitions.
//
// `EML` objects contain metadata about the dataset, such as the project name,
// URL, contact information, creators, and contributors.
type Importer interface {
	// Import takes a local path or a URL to a DwCA file (zip or tar.gz) and a
	// destination directory where the DwCA files will be extracted.
	//
	// srcPath: The local path or URL to the DwCA file.
	// dstDir: The destination directory for the extracted files.
	//
	// It downloads the file if a URL is provided and extracts it to the
	// destination directory.
	//
	// Returns: A directory in cache where meta.xml and eml.xml are located. Or,
	// it returns an error if the download, extraction, or initialization fails.
	Import(srcPath, dstDir string) (rootDir string, err error)

	// IngestMeta parses the meta.xml file found within a DwCA archive directory
	// and populates a *dwca.Meta object with its data.
	//
	// rootDir: The root directory of the extracted DwCA archive.
	//
	// Returns: A pointer to a populated *dwca.Meta object, or an error if the
	// meta.xml file is invalid, not found, or cannot be parsed.
	IngestMeta(rootDir string) (*Meta, error)

	// IngestMeta parses the eml.xml file found within a DwCA archive directory
	// and populates a *dwca.EML object with its data.
	//
	// rootDir: The root directory of the extracted DwCA archive.
	//
	// Returns: A pointer to a populated *dwca.EML object, or an error if the
	// eml.xml file is invalid, not found, or cannot be parsed.
	IngestEML(rootDir string) (*EML, error)
}

:w

type Accessor interface {
	Config() config.Config
}
