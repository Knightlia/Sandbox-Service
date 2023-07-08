package handlers

import (
	"net/http"

	"github.com/Knightlia/sandbox-service/model"
	"github.com/spf13/viper"
)

type HealthHandler struct{}

func NewHealthHandler() HealthHandler {
	return HealthHandler{}
}

func (_ HealthHandler) GetVersion(c model.Context) {
	c.PlainString(http.StatusOK, viper.GetString("version"))
}
