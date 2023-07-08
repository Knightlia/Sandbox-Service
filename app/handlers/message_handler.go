package handlers

import (
	"net/http"

	"github.com/Knightlia/sandbox-service/app/repository"
	"github.com/Knightlia/sandbox-service/cache"
	"github.com/Knightlia/sandbox-service/model"
	"github.com/go-chi/render"
)

type MessageHandler struct {
	userCache           cache.UserCache
	webSocketRepository repository.WebSocketRepository
}

func NewMessageHandler(userCache cache.UserCache, webSocketRepository repository.WebSocketRepository) MessageHandler {
	return MessageHandler{userCache, webSocketRepository}
}

func (m MessageHandler) SendMessage(c model.Context) {
	// Bind request model
	messageRequest := model.MessageRequest{}
	if err := render.Bind(c.Request(), &messageRequest); err != nil {
		c.JSON(http.StatusBadRequest, model.H{"error": err.Error()})
		return
	}

	// Get sender from token
	sender := m.userCache.Get(c.Request().Header.Get("token"))
	if sender == "" {
		c.JSON(http.StatusNotFound, model.H{"error": "error.sender.not.found"})
		return
	}

	// Broadcast message payload
	messagePayload := model.H{
		"messageType": "MESSAGE_PAYLOAD",
		"sender":      sender,
		"message":     messageRequest.Message,
		"timestamp":   messageRequest.Timestamp,
	}
	m.webSocketRepository.Broadcast(messagePayload)

	// Return response
	c.JSON(http.StatusOK, model.H{"status": "OK"})
}
