package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitFlags() {
	pflag.BoolP("version", "v", false, "Display application version.")
	pflag.UintP("port", "p", 8080, "The port to start the server on. Defauls to 8080.")
	pflag.Bool("debug", false, "Enable debug mode for development. False by default.")
	pflag.BoolP("help", "h", false, "Usage of Sandbox Service.")

	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
}

func HandleFlags(version string) bool {
	// Help flag
	if viper.GetBool("help") {
		pflag.Usage()
		return true
	}

	// Version flag
	if viper.GetBool("version") {
		fmt.Printf("Sandbox-Service : %s\n", version)
		return true
	}

	viper.Set("version", version)
	return false
}
