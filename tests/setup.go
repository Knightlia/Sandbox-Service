package tests

import (
	"net/http/httptest"

	"sandbox-service/app"
)

func SetupTests() *httptest.Server {
	a := app.NewApp()
	a.InitEcho()
	a.InitRoutes()

	return httptest.NewServer(a.Echo)
}
