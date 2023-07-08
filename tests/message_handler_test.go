package tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/Knightlia/sandbox-service/model"
	"github.com/stretchr/testify/assert"
)

func TestMessageHandler_SendMessage_401_WithMissingToken(t *testing.T) {
	s, _ := setup()
	defer s.Close()

	h := defaultHeaders()
	delete(h, "token")

	res := POST(s, "/message", h, messageRequest())

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assertJSONResponse(t, model.H{"error": "error.token.invalid"}, res.Body)

	teardown(res.Body)
}

func TestMessageHandler_SendMessage_401_WithInvalidToken(t *testing.T) {
	s, _ := setup()
	defer s.Close()

	h := defaultHeaders()
	h["token"] = "invalid-token"

	res := POST(s, "/message", h, messageRequest())

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assertJSONResponse(t, model.H{"error": "error.token.invalid"}, res.Body)

	teardown(res.Body)
}

func TestMessageHandler_SendMessage_400_WithMissingRequestBody(t *testing.T) {
	s, _ := setup()
	defer s.Close()

	res := POST(s, "/message", defaultHeaders(), nil)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assertJSONResponse(t, model.H{"error": "error.message.required"}, res.Body)

	teardown(res.Body)
}

func TestMessageHandler_SendMessage_400_WithInvalidMessage(t *testing.T) {
	s, _ := setup()
	defer s.Close()

	req := messageRequest()
	req.Message = ""

	res := POST(s, "/message", defaultHeaders(), req)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assertJSONResponse(t, model.H{"error": "error.message.required"}, res.Body)

	teardown(res.Body)
}

func TestMessageHandler_SendMessage_400_WithInvalidTimestamp(t *testing.T) {
	s, _ := setup()
	defer s.Close()

	req := messageRequest()
	req.Timestamp = 0

	res := POST(s, "/message", defaultHeaders(), req)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assertJSONResponse(t, model.H{"error": "error.timestamp.invalid"}, res.Body)

	teardown(res.Body)
}

func TestMessageHandler_SendMessage_404_IfUserDoesNotExist(t *testing.T) {
	s, _ := setup()
	defer s.Close()

	res := POST(s, "/message", defaultHeaders(), messageRequest())

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assertJSONResponse(t, model.H{"error": "error.sender.not.found"}, res.Body)

	teardown(res.Body)
}

func TestMessageHandler_SendMessage_200(t *testing.T) {
	s, a := setup()
	defer s.Close()
	a.UserCache.Store(validToken, "person-1")

	res := POST(s, "/message", defaultHeaders(), messageRequest())

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assertJSONResponse(t, model.H{"status": "OK"}, res.Body)

	teardown(res.Body)
}

func messageRequest() model.MessageRequest {
	return model.MessageRequest{
		Message:   "This is a test message.",
		Timestamp: time.Now().UnixMilli(),
	}
}
