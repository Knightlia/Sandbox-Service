package tests

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"sandbox-service/app"
)

func SetupTests() *echo.Echo {
	a := app.NewApp()
	a.InitEcho()
	a.InitRoutes()

	return a.Echo
}

func GET(e *echo.Echo, path string) *httptest.ResponseRecorder {
	return doRequest(e, http.MethodGet, path, nil)
}

func doRequest(e *echo.Echo, method string, path string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	res := httptest.NewRecorder()

	e.ServeHTTP(res, req)
	return res
}
