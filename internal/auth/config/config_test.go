package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	conf, err := Init()
	assert.NoError(t, err)
	assert.True(t, conf.Server.Port > 0)
	assert.True(t, len(conf.DBURI) > 0)
}
