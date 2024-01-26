package config

import (
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitLoggingLevel(t *testing.T) {
	cases := []struct {
		debug bool
		level zerolog.Level
	}{
		{false, zerolog.InfoLevel},
		{true, zerolog.DebugLevel},
	}

	for _, c := range cases {
		name := fmt.Sprintf("debug flag %t sets level %s", c.debug, c.level)
		t.Run(name, func(t *testing.T) {
			viper.Set("debug", c.debug)
			InitLogging()
			assert.Equal(t, c.level, zerolog.GlobalLevel())
			viper.Reset()
		})
	}
}
