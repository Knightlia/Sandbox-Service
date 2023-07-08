package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_GetVersion_200(t *testing.T) {
	s, _ := setup()
	defer s.Close()

	res := GET(s, "/")
	assert.Equal(t, http.StatusOK, res.StatusCode)

	teardown(res.Body)
}
