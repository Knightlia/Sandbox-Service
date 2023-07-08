package app

import (
	"net/http"

	"github.com/Knightlia/sandbox-service/app/handlers"
	"github.com/Knightlia/sandbox-service/app/repository"
	"github.com/Knightlia/sandbox-service/cache"
	"github.com/Knightlia/sandbox-service/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/olahol/melody"
	"github.com/spf13/viper"
)

type App struct {
	Chi    *chi.Mux
	Melody *melody.Melody

	UserCache cache.UserCache

	WebSocketRepository repository.WebSocketRepository
}

func NewApp(melody *melody.Melody) App {
	return App{
		Chi:    chi.NewRouter(),
		Melody: melody,

		UserCache: cache.NewUserCache(),

		WebSocketRepository: repository.NewWebSocketRepository(melody),
	}
}

// InitApp initialises the Chi router.
func (a App) InitApp() {
	// Global middleware
	a.Chi.Use(
		middleware.RequestID,
		// TODO: Logger with zerolog?
		middleware.Recoverer,

		cors.Handler(cors.Options{
			Debug:          viper.GetBool("debug"),
			AllowedOrigins: viper.GetStringSlice("cors"),
			AllowedHeaders: []string{"Content-Type", "token"},
		}),
	)
}

// InitRoutes initialises the handlers and the endpoints.
func (a App) InitRoutes() {
	healthHandler := handlers.NewHealthHandler()
	webSocketHandler := handlers.NewWebSocketHandler(a.Melody, a.UserCache, a.WebSocketRepository)
	nicknameHandler := handlers.NewNicknameHandler(a.UserCache, a.WebSocketRepository)
	messageHandler := handlers.NewMessageHandler(a.UserCache, a.WebSocketRepository)

	a.Chi.Get("/", a.handler(healthHandler.GetVersion))
	a.Chi.Get("/stream", a.handler(webSocketHandler.Connect))

	a.Chi.Group(func(r chi.Router) {
		r.Use(a.TokenMiddleware)
		r.Post("/nickname", a.handler(nicknameHandler.SetNickname))
		r.Post("/message", a.handler(messageHandler.SendMessage))
	})
}

// Wraps standard http handlers with [model.Context] used by handlers.
func (_ App) handler(h func(model.Context)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h(model.NewContext(w, r))
	}
}
