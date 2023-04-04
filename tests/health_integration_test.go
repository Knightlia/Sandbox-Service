package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthVersion(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	res, err := http.Get(s.URL)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
	}

	_ = res.Body.Close()
}
