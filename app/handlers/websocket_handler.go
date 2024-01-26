package handlers

import (
	"net/http"

	"github.com/olahol/melody"
	"github.com/rs/zerolog/log"
)

type WebSocketHandler struct {
	melody *melody.Melody
}

func NewWebSocketHandler(melody *melody.Melody) WebSocketHandler {
	w := WebSocketHandler{melody}
	melody.HandleConnect(w.handleConnect)
	melody.HandleDisconnect(w.handleDisconnect)
	return w
}

func (wh WebSocketHandler) Connect(w http.ResponseWriter, r *http.Request) {
	if err := wh.melody.HandleRequest(w, r); err != nil {
		log.Error().
			Err(err).
			Msg("Error handling websocket request.")
	}
}

func (_ WebSocketHandler) handleConnect(session *melody.Session) {
	log.Debug().Msg("New websocket client connected.")
}

func (_ WebSocketHandler) handleDisconnect(session *melody.Session) {
	log.Debug().Msg("Websocket client disconnected.")
}
