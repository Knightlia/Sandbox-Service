package repository

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WebSocketRepository struct {
	sessionRepository SessionRepository
}

func NewWebSocketRepository(sessionRepository SessionRepository) WebSocketRepository {
	return WebSocketRepository{sessionRepository}
}

func (w WebSocketRepository) SendPayload(ctx context.Context, conn *websocket.Conn, payload interface{}) error {
	c, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := wsjson.Write(c, conn, payload); err != nil {
		log.Error().Err(err).Msg("Failed to send websocket message.")
		return err
	}

	return nil
}

func (w WebSocketRepository) Broadcast(ctx context.Context, payload interface{}) {
	for _, conn := range w.sessionRepository.cache.Keys() {
		_ = w.SendPayload(ctx, conn, payload)
	}
}
