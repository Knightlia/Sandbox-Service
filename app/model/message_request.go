package model

import (
	"errors"
	"net/http"
)

type MessageRequest struct {
	Message   string
	Timestamp int64
}

func (m MessageRequest) Bind(_ *http.Request) error {
	if len(m.Message) == 0 {
		return errors.New("error.message.required")
	}
	if m.Timestamp <= 0 {
		return errors.New("error.timestamp.invalid")
	}
	return nil
}
