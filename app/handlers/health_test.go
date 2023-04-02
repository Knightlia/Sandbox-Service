package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"sandbox-service/tests"
)

func TestHealthHandler_Version(t *testing.T) {
	e := tests.SetupTests()
	res := tests.GET(e, "/")
	assert.Equal(t, http.StatusOK, res.Code)
}
