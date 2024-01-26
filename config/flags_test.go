package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitFlags(t *testing.T) {
	InitFlags()
	assert.Len(t, viper.AllKeys(), 4)

	reset()
}

func TestHandleFlagsHelp(t *testing.T) {
	InitFlags()
	viper.Set("help", true)

	assert.True(t, HandleFlags("test"))

	reset()
}

func TestHandleFlagsVersion(t *testing.T) {
	InitFlags()
	viper.Set("version", true)

	assert.True(t, HandleFlags("test"))

	reset()
}

func reset() {
	viper.Reset()
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
}
