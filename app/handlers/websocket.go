package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"sandbox-service/app/model"
	"sandbox-service/app/repository"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type WebSocketHandler struct {
	webSocketRepository repository.WebSocketRepository
	sessionRepository   repository.SessionRepository
	userRepository      repository.UserRepository
}

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
	u := websocket.NewUpgrader()
	u.CheckOrigin = w.checkOrigin
	u.KeepaliveTime = time.Second * 60

	u.OnOpen(w.onOpen)
	u.OnMessage(func(conn *websocket.Conn, messageType websocket.MessageType, bytes []byte) {
		message := string(bytes)
		if message == "ping" {
			_ = conn.SetReadDeadline(time.Now().Add(time.Second * 60))
		}
	})
	u.OnClose(w.onClose)

	_, err := u.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade websocket connection.")
		return
	}
}

func (w WebSocketHandler) checkOrigin(r *http.Request) bool {
	origins := viper.GetStringSlice("origins")
	o := r.Header.Get("Origin")

	if o != "" {
		for _, origin := range origins {
			if origin == o {
				return true
			}
		}
	}

	return false
}

func (w WebSocketHandler) onOpen(conn *websocket.Conn) {
	// Generate and store token
	token := generateSecureToken(32)
	w.sessionRepository.Store(conn, token)

	// Send token payload
	if err := w.webSocketRepository.SendPayload(conn, model.TokenPayload{
		MessageType: "TOKEN_PAYLOAD",
		Token:       token,
	}); err != nil {
		log.Error().Err(err).Msg("Failed to publish token payload.")
		_ = conn.Close()
	}

	// Send user list payload
	if err := w.webSocketRepository.SendPayload(conn, model.UserListPayload{
		MessageType: "USER_LIST_PAYLOAD",
		UserList:    w.userRepository.Values(),
	}); err != nil {
		log.Error().Err(err).Msg("Failed to publish user list payload snapshot.")
		_ = conn.Close()
	}
}

func (w WebSocketHandler) onClose(conn *websocket.Conn, err error) {
	if err != nil {
		log.Warn().Err(err).Msg("Websocket closed with error.")
	}
	log.Debug().Msg("Websocket connection closed.")
	v := w.sessionRepository.Remove(conn)
	w.userRepository.Remove(v)

	w.webSocketRepository.Broadcast(model.UserListPayload{
		MessageType: "USER_LIST_PAYLOAD",
		UserList:    w.userRepository.Values(),
	})
}

// Generates a random alphanumerical string with a specified length.
func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
