package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"sandbox-service/app"
	"sandbox-service/cache"
	"sandbox-service/config"
)

func main() {
	c := config.NewConfig()
	c.InitFlags()
	c.InitLogger()
	c.InitConfigFile()

	cache.InitCaches()

	a := app.NewApp()
	a.InitApp()
	a.InitRoutes()

	log.Fatal().
		Err(http.ListenAndServe(fmt.Sprintf(":%d", viper.GetUint("port")), a.Chi)).
		Msg("Failed to start server.")
}
