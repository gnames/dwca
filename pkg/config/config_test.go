package config_test

import (
	"testing"

	"github.com/gnames/dwca/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigDefault(t *testing.T) {
	assert := assert.New(t)
	conf := config.New()
	assert.Contains(conf.RootPath, "dwca_go")

	opts := []config.Option{
		config.OptRootPath("test"),
	}
	conf = config.New(opts...)
	assert.Equal("test", conf.RootPath)
}
