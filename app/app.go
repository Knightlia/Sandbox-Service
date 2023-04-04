package app

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/spf13/viper"
	"sandbox-service/app/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	Echo *echo.Echo
}

func NewApp() App {
	return App{Echo: echo.New()}
}

func (a App) InitEcho() {
	if viper.GetBool("debug") {
		a.Echo.Debug = true
	}

	a.Echo.Use(
		middleware.RequestID(),
		// TODO: middleware.Logger(),
		middleware.Recover(),
	)

	a.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
	}))
}

func (a App) InitRoutes() {
	healthHandler := handlers.NewHealthHandler()
    websocketHandler := handlers.NewWebSocketHandler()

	a.Echo.GET("/", healthHandler.Version)
	a.Echo.GET("/stream", websocketHandler.Connect)
}

func (a App) EnableMetrics() {
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(a.Echo)
}
