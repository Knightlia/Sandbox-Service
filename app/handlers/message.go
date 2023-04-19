package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"sandbox-service/app/model"
	"sandbox-service/app/repository"
)

type MessageHandler struct {
	userRepository      repository.UserRepository
	webSocketRepository repository.WebSocketRepository
}

func NewMessageHandler(userRepository repository.UserRepository, webSocketRepository repository.WebSocketRepository) MessageHandler {
	return MessageHandler{userRepository, webSocketRepository}
}

func (m MessageHandler) SendMessage(c model.Context) {
	// Bind to model request
	messageRequest := model.MessageRequest{}
	if err := render.Bind(c.Request(), &messageRequest); err != nil {
		c.JSON(http.StatusBadRequest, model.H{"error": err.Error()})
		return
	}

	// Get sender from token
	sender := m.userRepository.Get(c.Request().Header.Get("token"))
	if sender == "" {
		c.JSON(http.StatusNotFound, model.H{"error": "error.sender.not.found"})
		return
	}

	// Broadcast message
	m.webSocketRepository.Broadcast(c.Request().Context(), model.MessagePayload{
		MessageType: "MESSAGE_PAYLOAD",
		Sender:      sender,
		Message:     messageRequest.Message,
		Timestamp:   messageRequest.Timestamp,
	})

	// Return response
	c.JSON(http.StatusOK, model.H{"status": "OK"})
}
