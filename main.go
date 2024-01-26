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
	"github.com/olahol/melody"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var version = "dev"

func main() {
	// Initialise configuration
	config.InitFlags()
	if config.HandleFlags(version) {
		return
	}
	config.InitLogging()

	// Initialise melody
	melody := melody.New()
	defer melody.Close()

	// Initialise server and routes
	app := app.NewApp(melody)
	mux := app.Init()
	app.Routes(mux)

	// Start server
	addr := fmt.Sprintf(":%d", viper.GetUint("port"))
	server := http.Server{Addr: addr, Handler: mux}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panic().
				Err(err).
				Msg("Failed to start server.")
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
		log.Panic().
			Err(err).
			Msg("Error during graceful shutdown.")
	}
	log.Info().Msg("Server stopped.")
}
