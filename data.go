package dwca

import "github.com/gnames/dwca/ent/eml"

type ArchiveData struct {
	eml.EML
	Meta
	Core
	Extensions []Extension
}

type Meta struct {
}

type Core struct {
}

type Extension struct {
}
