package impio

import "github.com/gnames/dwca"

type impio struct{}

func New() dwca.Importer {
	return &impio{}
}

func (i *impio) IngestMeta(rootDir string) (*dwca.Meta, error) {
	return nil, nil
}

func (i *impio) IngestEML(rootDir string) (*dwca.EML, error) {
	return nil, nil
}
