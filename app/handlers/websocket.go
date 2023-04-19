package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"sandbox-service/app/model"
	"sandbox-service/app/repository"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WebSocketHandler struct {
	webSocketRepository repository.WebSocketRepository
	sessionRepository   repository.SessionRepository
	userRepository      repository.UserRepository
}

// NewWebSocketHandler creates a new instance of the websocket handler.
func NewWebSocketHandler(
	webSocketRepository repository.WebSocketRepository,
	sessionRepository repository.SessionRepository,
	userRepository repository.UserRepository,
) WebSocketHandler {
	return WebSocketHandler{webSocketRepository, sessionRepository, userRepository}
}

// Connect accepts incoming websocket connections from clients. Takes the
// [model.Context] as the parameter.
func (w WebSocketHandler) Connect(c model.Context) {
	// Accept websocket connection
	conn, err := websocket.Accept(c.Response(), c.Request(), &websocket.AcceptOptions{
		OriginPatterns: viper.GetStringSlice("origins.ws"),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to websocket.")
		return
	}
	defer w.close(conn)

	// Store session
	token := generateSecureToken(32)
	w.sessionRepository.Store(conn, token)

	// Send token payload
	if err := w.sendTokenPayload(c, conn, token); err != nil {
		log.Error().Err(err).Msg("Failed to publish token payload.")
		return
	}

	// Send user list payload
	w.sendUserListPayload(c, conn)

	// Loop to keep the connection alive
	for {
		var v interface{}
		_ = wsjson.Read(context.Background(), conn, &v)
	}
}

// Attempts to cleanly close a websocket connection and handle any errors.
func (w WebSocketHandler) close(conn *websocket.Conn) {
	log.Debug().Msg("Closing websocket connection...")
	if err := conn.Close(websocket.StatusNormalClosure, ""); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to close websocket connection.")
	}

	v := w.sessionRepository.Remove(conn)
	w.userRepository.Remove(v)
}

// The sendTokenPayload method attempts to publish a [model.TokenPayload] to the
// client connecting for the first time.
func (w WebSocketHandler) sendTokenPayload(c model.Context, conn *websocket.Conn, token string) error {
	err := w.webSocketRepository.SendPayload(c.Request().Context(), conn, model.TokenPayload{
		MessageType: "TOKEN_PAYLOAD",
		Token:       token,
	})
	return err
}

// The sendUserListPayload publishes a [model.UserListPayload] which contains a snapshot
// of all the nicknames.
func (w WebSocketHandler) sendUserListPayload(c model.Context, conn *websocket.Conn) {
	_ = w.webSocketRepository.SendPayload(c.Request().Context(), conn, model.UserListPayload{
		MessageType: "USER_LIST_PAYLOAD",
		UserList:    w.userRepository.Values(),
	})
}

// ---

// Generates a random alphanumerical string with a specified length.
func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
