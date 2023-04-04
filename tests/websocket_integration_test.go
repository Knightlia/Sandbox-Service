package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
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

		var payload map[string]interface{}
		err = wsjson.Read(ctx, c, &payload)

		if assert.NoError(t, err) {
			fmt.Println(payload)
			assert.Equal(t, "TOKEN_PAYLOAD", payload["messageType"])
		}
	}
}

func closeConn(c *websocket.Conn) {
	_ = c.Close(websocket.StatusNormalClosure, "")
}
