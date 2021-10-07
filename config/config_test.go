package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_ReadConfig(t *testing.T) {
	cfg := &Config{}
	err := cfg.ReadConfig("./config.example.yaml")
	assert.NoError(t, err)
	t.Log(cfg)
}
