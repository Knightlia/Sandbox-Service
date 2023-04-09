package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"sandbox-service/app/model"
	"sandbox-service/cache"
)

func TestMessage401WithMissingToken(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	headers := DefaultHeaders()
	delete(headers, "token")

	res := POST(s, "/message", headers, messageRequest())
	defer TeardownTests(res.Body)

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	AssertJSONResponse(t, model.H{"error": "error.token.invalid"}, res.Body)
}

func TestMessage401WithInvalidToken(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	headers := DefaultHeaders()
	headers["token"] = "invalid-token"

	res := POST(s, "/message", headers, nil)
	defer TeardownTests(res.Body)

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	AssertJSONResponse(t, model.H{"error": "error.token.invalid"}, res.Body)
}

func TestMessage400WithNoRequest(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	res := POST(s, "/message", DefaultHeaders(), nil)
	defer TeardownTests(res.Body)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	AssertJSONResponse(t, model.H{"error": "error.message.required"}, res.Body)
}

func TestMessage400WithInvalidRequestMessage(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	req := messageRequest()
	req.Message = ""

	res := POST(s, "/message", DefaultHeaders(), req)
	defer TeardownTests(res.Body)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	AssertJSONResponse(t, model.H{"error": "error.message.required"}, res.Body)
}

func TestMessage400WithInvalidRequestTimestamp(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	req := messageRequest()
	req.Timestamp = 0

	res := POST(s, "/message", DefaultHeaders(), req)
	defer TeardownTests(res.Body)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	AssertJSONResponse(t, model.H{"error": "error.timestamp.invalid"}, res.Body)
}

func TestMessage404IfSenderNotFound(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	cache.UserCache.DeleteAll()

	res := POST(s, "/message", DefaultHeaders(), messageRequest())
	defer TeardownTests(res.Body)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	AssertJSONResponse(t, model.H{"error": "error.sender.not.found"}, res.Body)
}

func TestMessage200(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	res := POST(s, "/message", DefaultHeaders(), messageRequest())
	defer TeardownTests(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	AssertJSONResponse(t, model.H{"status": "OK"}, res.Body)
}

func messageRequest() model.MessageRequest {
	return model.MessageRequest{
		Message:   "A test message.",
		Timestamp: 1680991825224,
	}
}
