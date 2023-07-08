package app

import (
	"net/http"

	"github.com/Knightlia/sandbox-service/model"
	"github.com/go-chi/render"
)

func (a App) TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")

		if len(token) == 0 || !a.UserCache.HasKey(token) {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, model.H{"error": "error.token.invalid"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
