package dwca

import (
	"github.com/gnames/dwca/config"
	"github.com/gnames/dwca/ent/meta"
)

// Archive is an interface for Darwin Core Archive objects.
type Archive interface {
	Load() error
	Meta() *meta.Meta
	Config() config.Config
}
