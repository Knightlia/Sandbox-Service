package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"sandbox-service/app"
	"sandbox-service/config"
)

func main() {
	c := config.NewConfig()
	c.InitFlags()
	c.InitLogger()
	c.InitConfigFile()

	a := app.NewApp()
	a.InitEcho()
	a.InitRoutes()
	a.EnableMetrics()

	log.Fatal().
		Err(a.Echo.Start(fmt.Sprintf(":%d", viper.GetUint("port")))).
		Msg("Failed to start server.")
}
