package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const version = "0.0.1"

type HealthHandler struct{}

func NewHealthHandler() HealthHandler {
	return HealthHandler{}
}

func (_ HealthHandler) Version(c echo.Context) error {
	return c.String(http.StatusOK, version)
}
