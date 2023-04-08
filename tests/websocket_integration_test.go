package tests

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"sandbox-service/app/model"
)

func TestWebSocketConnect(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	url := strings.Replace(s.URL+"/stream", "http", "ws", 1)
	c, _, err := websocket.Dial(context.Background(), url, nil)
	defer closeConn(c)

	if assert.NoError(t, err) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Check for token payload
		payload := make(map[string]interface{}, 0)
		err = wsjson.Read(ctx, c, &payload)
		if assert.NoError(t, err) {
			assert.Equal(t, "TOKEN_PAYLOAD", payload["mt"])
		}

		// Check for user List payload
		payload = make(map[string]interface{}, 0)
		err = wsjson.Read(ctx, c, &payload)
		if assert.NoError(t, err) {
			assert.Equal(t, "USER_LIST_PAYLOAD", payload["mt"])
		}
	}
}

func TestWebSocketNickname(t *testing.T) {
	s := SetupTests()
	defer s.Close()

	url := strings.Replace(s.URL+"/stream", "http", "ws", 1)
	c, _, err := websocket.Dial(context.Background(), url, nil)
	defer closeConn(c)

	if assert.NoError(t, err) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Get token
		payload := make(map[string]interface{}, 0)
		_ = wsjson.Read(ctx, c, &payload)
		token := payload["t"]

		// Drop initial user list payloads
		_ = wsjson.Read(ctx, c, &model.H{})

		// Send nickname request
		headers := DefaultHeaders()
		headers["token"] = token.(string)
		res := POST(s, "/nickname", headers, nicknameRequest())
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// Consume user list payload
		payload = make(map[string]interface{}, 0)
		err := wsjson.Read(ctx, c, &payload)
		if assert.NoError(t, err) {
			assert.Equal(t, "USER_LIST_PAYLOAD", payload["mt"])
			assert.ElementsMatch(t, []string{"nickname", "existing-nickname"}, payload["ul"])
		}
	}
}

func closeConn(c *websocket.Conn) {
	_ = c.Close(websocket.StatusNormalClosure, "")
}
