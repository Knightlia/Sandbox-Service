package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Knightlia/sandbox-service/app"
	"github.com/Knightlia/sandbox-service/config"
	"github.com/getsentry/sentry-go"
	"github.com/olahol/melody"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	version   string
	commit    string
	buildTime string
)

func main() {
	initConfiguration()

	// Initialise [melody] for websockets.
	m := melody.New()
	defer func() {
		if err := m.Close(); err != nil {
			sentry.CaptureException(err)
			log.Error().
				Err(err).
				Msg("Error while closing the Melody (websocket) instance.")
		}
	}()

	// Setup application
	a := app.NewApp(m)
	a.InitApp()
	a.InitRoutes()

	// Initialise Sentry
	if viper.GetBool("sentry.enabled") {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              viper.GetString("sentry.dsn"),
			Debug:            viper.GetBool("debug"),
			AttachStacktrace: true,
			EnableTracing:    true,
			TracesSampleRate: viper.GetFloat64("sentry.traces_sample_rate"),
			Environment:      viper.GetString("environment"),
		}); err != nil {
			log.Panic().
				Err(err).
				Msg("Error while initialising Sentry.")
		}

		defer sentry.Flush(2 * time.Second)
	}

	// Start server
	addr := fmt.Sprintf(":%d", viper.GetUint("port"))
	server := http.Server{Addr: addr, Handler: a.Chi}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sentry.CaptureException(err)
			log.Panic().
				Err(err).
				Msg("Shutting down server.")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info().Msg("Performing graceful shutdown...")
	if err := server.Shutdown(ctx); err != nil {
		sentry.CaptureException(err)
		log.Panic().
			Err(err).
			Msg("Error during graceful shutdown.")
	}
	log.Info().Msg("Server stopped.")
}

// Initialises application CLI flags and loads configurations.
func initConfiguration() {
	c := config.NewConfig()
	c.InitFlags()

	if viper.GetBool("version") || viper.GetBool("v") {
		fmt.Println("\n=== Sandbox Service ===")
		fmt.Printf("Version:     %s\n", version)
		fmt.Printf("Commit:      %s\n", commit)
		fmt.Printf("Build Time:  %s\n", buildTime)
		return
	}

	viper.Set("version", version)

	c.InitLogger()
	c.InitConfigFile()

	// Print config (debug)
	fmt.Println()
	for _, v := range viper.AllKeys() {
		log.Debug().Msgf("%s: %s", v, viper.Get(v))
	}
	fmt.Println()
}
