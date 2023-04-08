package middleware

import (
	"net/http"

	"github.com/go-chi/render"
	"sandbox-service/app/model"
	"sandbox-service/app/repository"
)

type Middleware struct {
	sessionRepository repository.SessionRepository
}

func NewMiddleware(sessionRepository repository.SessionRepository) Middleware {
	return Middleware{sessionRepository}
}

func (m Middleware) TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")

		if len(token) == 0 || !m.sessionRepository.HasValue(token) {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, model.H{"error": "error.token.invalid"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
