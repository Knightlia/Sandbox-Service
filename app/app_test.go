package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestApp_InitEcho(t *testing.T) {
	viper.Set("debug", true)

	a := NewApp()
	a.InitEcho()

	assert.True(t, a.Echo.Debug)
}

func TestApp_InitRoutes(t *testing.T) {
	a := NewApp()
	a.InitRoutes()

	assert.GreaterOrEqual(t, len(a.Echo.Routes()), 1)
}

func TestApp_EnableMetrics(t *testing.T) {
	a := NewApp()
	a.EnableMetrics()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	res := httptest.NewRecorder()
	a.Echo.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}
