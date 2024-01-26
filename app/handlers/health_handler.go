package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/spf13/viper"
)

type HealthHandler struct{}

func NewHealthHandler() HealthHandler {
	return HealthHandler{}
}

func (_ HealthHandler) GetVersion(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, viper.GetString("version"))
}
