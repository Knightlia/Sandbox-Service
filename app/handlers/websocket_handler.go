package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/Knightlia/sandbox-service/app/repository"
	"github.com/Knightlia/sandbox-service/cache"
	"github.com/Knightlia/sandbox-service/model"
	"github.com/getsentry/sentry-go"
	"github.com/olahol/melody"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

type WebSocketHandler struct {
	melody              *melody.Melody
	userCache           cache.UserCache
	webSocketRepository repository.WebSocketRepository
}

func NewWebSocketHandler(
	melody *melody.Melody,
	userCache cache.UserCache,
	webSocketRepository repository.WebSocketRepository,
) WebSocketHandler {
	w := WebSocketHandler{melody, userCache, webSocketRepository}

	melody.Upgrader.CheckOrigin = w.checkOrigin
	melody.HandleConnect(w.onConnect)
	melody.HandleDisconnect(w.onDisconnect)

	return w
}

func (w WebSocketHandler) Connect(c model.Context) {
	if err := w.melody.HandleRequest(c.Response(), c.Request()); err != nil {
		sentry.CaptureException(err)
		log.Error().
			Err(err).
			Msg("Error handling websocket request.")
	}
}

// Runs when a new websocket client connects.
func (w WebSocketHandler) onConnect(session *melody.Session) {
	log.Debug().Msg("New websocket client connected.")

	token := w.generateSecureToken()
	tokenPayload := model.H{
		"messageType": "TOKEN_PAYLOAD",
		"token":       token,
	}

	if err := w.webSocketRepository.SendSinglePayload(session, tokenPayload); err != nil {
		w.closeSession(session)
		return
	}

	w.broadcastUserList()

	session.Set("token", token)
	w.userCache.Store(token, "")
}

// Runs when a websocket client connection closes.
func (w WebSocketHandler) onDisconnect(session *melody.Session) {
	log.Debug().Msg("Websocket client disconnected.")

	token, exists := session.Get("token")
	if exists {
		log.Debug().Msgf("Removing token %s from cache.", token)
		w.userCache.Remove(token.(string))
		w.broadcastUserList()
	}
}

// Helper function to generate a random 64 digit token.
func (_ WebSocketHandler) generateSecureToken() string {
	b := make([]byte, 64)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Helper function to close websocket sessions if they are not closed already.
func (_ WebSocketHandler) closeSession(session *melody.Session) {
	if err := session.Close(); err != nil {
		sentry.CaptureException(err)
		log.Error().
			Err(err).
			Msg("Error closing websocket session.")
	}
}

func (_ WebSocketHandler) checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("origin")
	if len(origin) == 0 {
		return true
	}

	return slices.Contains(viper.GetStringSlice("cors"), origin)
}

// Helper function to create a user list payload and broadcast to all websocket clients.
func (w WebSocketHandler) broadcastUserList() {
	userListPayload := model.H{
		"messageType": "USER_LIST_PAYLOAD",
		"userList":    w.userCache.Values(),
	}
	w.webSocketRepository.Broadcast(userListPayload)
}
