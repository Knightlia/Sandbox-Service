package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"sandbox-service/app/model"
	"sandbox-service/app/repository"
)

type NicknameHandler struct {
	userRepository      repository.UserRepository
	webSocketRepository repository.WebSocketRepository
}

func NewNicknameHandler(userRepository repository.UserRepository, webSocketRepository repository.WebSocketRepository) NicknameHandler {
	return NicknameHandler{userRepository, webSocketRepository}
}

func (n NicknameHandler) SetNickname(c model.Context) {
	// Bind request model
	nicknameRequest := model.NicknameRequest{}
	if err := render.Bind(c.Request(), &nicknameRequest); err != nil {
		c.JSON(http.StatusBadRequest, model.H{"error": err.Error()})
		return
	}

	// Check if nickname exists
	if n.userRepository.HasValue(nicknameRequest.Nickname) {
		c.JSON(http.StatusOK, model.H{"error": "error.nickname.exists"})
		return
	}

	// Store nickname
	n.userRepository.Store(c.Request().Header.Get("token"), nicknameRequest.Nickname)

	// Broadcast userList
	n.webSocketRepository.Broadcast(model.UserListPayload{
		MessageType: "USER_LIST_PAYLOAD",
		UserList:    n.userRepository.Values(),
	})

	// Return response
	c.JSON(http.StatusOK, model.H{"status": "OK"})
}
