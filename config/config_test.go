package config

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfig_InitFlags(t *testing.T) {
	c := NewConfig()
	c.InitFlags()

	// Test one of the flags
	assert.Equal(t, 8080, viper.GetInt("port"))
}

func TestConfig_InitLogger(t *testing.T) {
	c := NewConfig()
	c.InitLogger()
	assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())

	viper.Set("debug", true)
	c.InitLogger()
	assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
}

func TestConfig_InitConfigFile(t *testing.T) {
	c := NewConfig()
	assert.Panics(t, func() {
		c.InitConfigFile()
	})

	viper.AddConfigPath("../")
	c.InitConfigFile()
	assert.NotEmpty(t, viper.ConfigFileUsed())
}
