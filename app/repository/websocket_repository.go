package repository

import (
	"encoding/json"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/rs/zerolog/log"
)

type WebSocketRepository struct {
	sessionRepository SessionRepository
}

func NewWebSocketRepository(sessionRepository SessionRepository) WebSocketRepository {
	return WebSocketRepository{sessionRepository}
}

func (w WebSocketRepository) SendPayload(conn *websocket.Conn, payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal websocket message.")
		return err
	}

	return w.sendPayload(conn, b)
}

func (w WebSocketRepository) Broadcast(payload interface{}) {
	b, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal websocket message.")
		return
	}

	for _, conn := range w.sessionRepository.cache.Keys() {
		_ = w.sendPayload(conn, b)
	}
}

func (w WebSocketRepository) sendPayload(conn *websocket.Conn, payload []byte) error {
	if err := conn.WriteMessage(websocket.BinaryMessage, payload); err != nil {
		log.Error().Err(err).Msg("Failed to send websocket message.")
		return err
	}

	return nil
}
