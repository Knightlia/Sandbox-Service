package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"sandbox-service/app/model"
)

func TestNickname401WithMissingToken(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	headers := DefaultHeaders()
	delete(headers, "token")

	res := POST(s, "/nickname", headers, nil)
	defer TeardownTests(res.Body)

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	AssertJSONResponse(t, model.H{"error": "error.token.invalid"}, res.Body)
}

func TestNickname401WithInvalidToken(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	headers := DefaultHeaders()
	headers["token"] = "invalid-token"

	res := POST(s, "/nickname", headers, nicknameRequest())
	defer TeardownTests(res.Body)

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	AssertJSONResponse(t, model.H{"error": "error.token.invalid"}, res.Body)
}

func TestNickname400WithNoRequest(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	res := POST(s, "/nickname", DefaultHeaders(), nil)
	defer TeardownTests(res.Body)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	AssertJSONResponse(t, model.H{"error": "error.nickname.required"}, res.Body)
}

func TestNickname400WithInvalidRequest(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	req := nicknameRequest()
	req.Nickname = "invalid nickname"

	res := POST(s, "/nickname", DefaultHeaders(), req)
	defer TeardownTests(res.Body)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	AssertJSONResponse(t, model.H{"error": "error.nickname.format"}, res.Body)
}

func TestNickname200WithExistingUserName(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	req := nicknameRequest()
	req.Nickname = "existing-nickname"

	res := POST(s, "/nickname", DefaultHeaders(), req)
	defer TeardownTests(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	AssertJSONResponse(t, model.H{"error": "error.nickname.exists"}, res.Body)
}

func TestNickname200(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	res := POST(s, "/nickname", DefaultHeaders(), nicknameRequest())
	defer TeardownTests(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	AssertJSONResponse(t, model.H{"status": "OK"}, res.Body)
}

func nicknameRequest() model.NicknameRequest {
	return model.NicknameRequest{Nickname: "nickname"}
}
