package tests

import (
	"net/http"
	"testing"

	"github.com/Knightlia/sandbox-service/model"
	"github.com/stretchr/testify/assert"
)

func TestNicknameHandler_SetNickname_401_WithMissingToken(t *testing.T) {
	s, _ := setup()
	defer s.Close()

	h := defaultHeaders()
	delete(h, "token")

	res := POST(s, "/nickname", h, nicknameRequest())

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assertJSONResponse(t, model.H{"error": "error.token.invalid"}, res.Body)

	teardown(res.Body)
}

func TestNicknameHandler_SetNickname_401_WithInvalidToken(t *testing.T) {
	s, _ := setup()
	defer s.Close()

	h := defaultHeaders()
	h["token"] = "invalid-token"

	res := POST(s, "/nickname", h, nicknameRequest())

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assertJSONResponse(t, model.H{"error": "error.token.invalid"}, res.Body)

	teardown(res.Body)
}

func TestNicknameHandler_SetNickname_400_WithMissingRequestBody(t *testing.T) {
	s, _ := setup()
	defer s.Close()

	res := POST(s, "/nickname", defaultHeaders(), nil)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assertJSONResponse(t, model.H{"error": "error.nickname.required"}, res.Body)

	teardown(res.Body)
}

func TestNicknameHandler_SetNickname_400_WithInvalidRequestBody(t *testing.T) {
	s, _ := setup()
	defer s.Close()

	req := nicknameRequest()
	req.Nickname = "invalid nickname"

	res := POST(s, "/nickname", defaultHeaders(), req)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assertJSONResponse(t, model.H{"error": "error.nickname.format"}, res.Body)

	teardown(res.Body)
}

func TestNicknameHandler_SetNickname_400_IfNicknameExists(t *testing.T) {
	s, a := setup()
	defer s.Close()

	// Add user in the cache
	a.UserCache.Store("token", "person-1")

	res := POST(s, "/nickname", defaultHeaders(), nicknameRequest())

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assertJSONResponse(t, model.H{"error": "error.nickname.exists"}, res.Body)

	teardown(res.Body)
}

func TestNicknameHandler_SetNickname_200(t *testing.T) {
	s, a := setup()
	defer s.Close()

	res := POST(s, "/nickname", defaultHeaders(), nicknameRequest())

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assertJSONResponse(t, model.H{"status": "OK"}, res.Body)
	assert.True(t, a.UserCache.HasValue("person-1"))

	teardown(res.Body)
}

func TestNicknameHandler_SetNickname_200_EvenIfBroadcastFails(t *testing.T) {
	s, a := setup()
	defer s.Close()

	_ = a.Melody.Close()

	res := POST(s, "/nickname", defaultHeaders(), nicknameRequest())

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assertJSONResponse(t, model.H{"status": "OK"}, res.Body)
	assert.True(t, a.UserCache.HasValue("person-1"))

	teardown(res.Body)
}

func nicknameRequest() model.NicknameRequest {
	return model.NicknameRequest{Nickname: "person-1"}
}
