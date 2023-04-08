package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"sandbox-service/app/handlers"
	m "sandbox-service/app/middleware"
	"sandbox-service/app/model"
	"sandbox-service/app/repository"
	"sandbox-service/cache"
)

var (
	WebSocketRepository repository.WebSocketRepository
	SessionRepository   repository.SessionRepository
	UserRepository      repository.UserRepository
)

type App struct {
	Chi        *chi.Mux
	middleware m.Middleware
}

func NewApp() App {
	// Initialise repositories
	SessionRepository = repository.NewSessionRepository(cache.SessionCache)
	UserRepository = repository.NewUserRepository(cache.UserCache)
	WebSocketRepository = repository.NewWebSocketRepository(SessionRepository)

	return App{Chi: chi.NewRouter(), middleware: m.NewMiddleware(SessionRepository)}
}

// InitApp initialises the Chi router.
func (a App) InitApp() {
	// Global middleware
	a.Chi.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
	)
}

func (a App) InitRoutes() {
	healthHandler := handlers.NewHealthHandler()
	webSocketHandler := handlers.NewWebSocketHandler(WebSocketRepository, SessionRepository, UserRepository)
	nicknameHandler := handlers.NewNicknameHandler(UserRepository, WebSocketRepository)

	a.Chi.Get("/", a.handler(healthHandler.Version))
	a.Chi.Get("/stream", a.handler(webSocketHandler.Connect))

	a.Chi.Group(func(r chi.Router) {
		r.Use(a.middleware.TokenMiddleware)
		r.Post("/nickname", a.handler(nicknameHandler.SetNickname))
	})
}

// ---

func (_ App) handler(h func(model.Context)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h(model.NewContext(w, r))
	}
}
