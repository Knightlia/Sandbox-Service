package tests

import (
	"strings"
	"testing"

	"github.com/Knightlia/sandbox-service/model"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketHandler_Connect_SendsTokenAndUserListOnConnect(t *testing.T) {
	s, a := setup()

	url := strings.Replace(s.URL+"/stream", "http", "ws", 1)
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(url, nil)

	if assert.NoError(t, err) {
		v := model.H{}
		if assert.NoError(t, conn.ReadJSON(&v)) {
			assert.Equal(t, "TOKEN_PAYLOAD", v["messageType"])
			assert.True(t, a.UserCache.HasKey(v["token"].(string)))
		}

		v = model.H{}
		if assert.NoError(t, conn.ReadJSON(&v)) {
			assert.Equal(t, "USER_LIST_PAYLOAD", v["messageType"])
		}
	}
}

func TestWebSocketHandler_Connect_FailsOnMelodyError(t *testing.T) {
	s, a := setup()
	_ = a.Melody.Close()

	url := strings.Replace(s.URL+"/stream", "http", "ws", 1)
	dialer := websocket.Dialer{}
	_, _, err := dialer.Dial(url, nil)

	assert.Error(t, err)
}

func TestWebSocketHandler_onDisconnect_RemovesUserFromCache(t *testing.T) {
	s, a := setup()

	url := strings.Replace(s.URL+"/stream", "http", "ws", 1)
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(url, nil)

	if assert.NoError(t, err) {
		assert.NoError(t, conn.ReadJSON(&model.H{}))
	}

	if assert.NoError(t, conn.Close()) {
		assert.Len(t, a.UserCache.Values(), 0)
	}
}
