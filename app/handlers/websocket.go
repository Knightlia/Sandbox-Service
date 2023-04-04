package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"sandbox-service/app/model"
)

type WebSocketHandler struct{}

// NewWebSocketHandler creates a new instance of the websocket handler.
func NewWebSocketHandler() WebSocketHandler {
	return WebSocketHandler{}
}

// Connect accepts incoming websocket connections from clients. Takes the
// [echo.Context] as the parameter.
func (wh WebSocketHandler) Connect(c echo.Context) error {
	// Accept websocket connection
	conn, err := websocket.Accept(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to websocket.")
		return err
	}
	defer wh.close(conn)

	// Send token payload
	if err := wh.sendTokenPayload(c, conn); err != nil {
		log.Error().Err(err).Msg("Failed to publish token payload.")
		return err
	}

	// Loop to keep the connection alive
	for {
		var v interface{}
		_ = wsjson.Read(context.Background(), conn, &v)
	}
}

// Attempts to cleanly close a websocket connection and handle any errors.
func (_ WebSocketHandler) close(conn *websocket.Conn) {
	if err := conn.Close(websocket.StatusNormalClosure, ""); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to close websocket connection.")
	}
}

// The sendTokenPayload method attempts to publish a [model.TokenPayload] to the
// client connecting for the first time.
func (_ WebSocketHandler) sendTokenPayload(c echo.Context, conn *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	err := wsjson.Write(ctx, conn, model.TokenPayload{
		MessageType: "TOKEN_PAYLOAD",
		Token:       generateSecureToken(64),
	})

	return err
}

// Generates a random alphanumerical string with a specified length.
func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
