package app

import (
	"github.com/Knightlia/sandbox-service/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/olahol/melody"
	"github.com/rs/zerolog/log"
)

type App struct {
	melody *melody.Melody
}

func NewApp(m *melody.Melody) App {
	return App{m}
}

func (_ App) Init() *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
	)

	return r
}

func (a App) Routes(r *chi.Mux) {
	healthHandler := handlers.NewHealthHandler()
	webSocketHandler := handlers.NewWebSocketHandler(a.melody)

	r.Get("/", healthHandler.GetVersion)
	r.Get("/stream", webSocketHandler.Connect)

	log.Info().Msgf("%d routes initialised.", len(r.Routes()))
}
