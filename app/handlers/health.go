package handlers

import (
	"net/http"

	"sandbox-service/app/model"
)

const version = "0.0.1"

type HealthHandler struct{}

func NewHealthHandler() HealthHandler {
	return HealthHandler{}
}

func (_ HealthHandler) Version(c model.Context) {
	c.PlainString(http.StatusOK, version)
}
