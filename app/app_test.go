package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitGlobalMiddleware(t *testing.T) {
	app := NewApp()
	r := app.Init()

	assert.GreaterOrEqual(t, len(r.Middlewares()), 1)
}

func TestRoutesInitialises(t *testing.T) {
	app := NewApp()
	r := app.Init()
	app.Routes(r)

	assert.GreaterOrEqual(t, len(r.Routes()), 1)
}
