package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Knightlia/sandbox-service/app"
	"github.com/stretchr/testify/assert"
)

func TestGetVersion200(t *testing.T) {
	app := app.NewApp()
	r := app.Init()
	app.Routes(r)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}
