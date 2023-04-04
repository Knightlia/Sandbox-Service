package tests

import (
	"testing"

	"sandbox-service/config"

	"github.com/stretchr/testify/assert"
)

func TestAppPanicsIfConfigFileNotFound(t *testing.T) {
	c := config.NewConfig()
	assert.Panics(t, func() {
		c.InitConfigFile()
	})
}
