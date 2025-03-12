package dwcio_test

import (
	"testing"

	"github.com/gnames/dwca/config"
	"github.com/gnames/dwca/io/dwcio"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)
	cfg := config.New()
	dwca := dwcio.New(cfg)
	cfg2 := dwca.Config()
	assert.Equal(cfg, cfg2)
}
