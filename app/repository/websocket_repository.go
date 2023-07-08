package repository

import (
	"encoding/json"

	"github.com/getsentry/sentry-go"
	"github.com/olahol/melody"
	"github.com/rs/zerolog/log"
)

type WebSocketRepository struct {
	melody *melody.Melody
}

func NewWebSocketRepository(melody *melody.Melody) WebSocketRepository {
	return WebSocketRepository{melody}
}

func (w WebSocketRepository) SendSinglePayload(session *melody.Session, payload interface{}) error {
	b, _ := json.Marshal(payload)

	if err := session.Write(b); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to send single websocket payload to session.")
		sentry.CaptureException(err)
		return err
	}

	return nil
}

func (w WebSocketRepository) Broadcast(payload interface{}) {
	b, _ := json.Marshal(payload)

	if err := w.melody.Broadcast(b); err != nil {
		sentry.CaptureException(err)
		log.Error().
			Err(err).
			Msg("Failed to broadcast websocket message.")
	}
}
