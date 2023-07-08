package handlers

import (
	"net/http"

	"github.com/Knightlia/sandbox-service/app/repository"
	"github.com/Knightlia/sandbox-service/cache"
	"github.com/Knightlia/sandbox-service/model"
	"github.com/go-chi/render"
)

type NicknameHandler struct {
	userCache           cache.UserCache
	webSocketRepository repository.WebSocketRepository
}

func NewNicknameHandler(userCache cache.UserCache, webSocketRepository repository.WebSocketRepository) NicknameHandler {
	return NicknameHandler{userCache, webSocketRepository}
}

func (n NicknameHandler) SetNickname(c model.Context) {
	// Bind and validate request
	nicknameRequest := model.NicknameRequest{}
	if err := render.Bind(c.Request(), &nicknameRequest); err != nil {
		c.JSON(http.StatusBadRequest, model.H{"error": err.Error()})
		return
	}

	// Check if name exists already
	if n.userCache.HasValue(nicknameRequest.Nickname) {
		c.JSON(http.StatusBadRequest, model.H{"error": "error.nickname.exists"})
		return
	}

	// Store nickname and broadcast
	n.userCache.Store(c.Header("token"), nicknameRequest.Nickname)
	n.broadcastUserList()

	// Return response
	c.JSON(http.StatusOK, model.H{"status": "OK"})
}

// Helper function to create a user list payload and broadcast to all websocket clients.
func (n NicknameHandler) broadcastUserList() {
	userListPayload := model.H{
		"messageType": "USER_LIST_PAYLOAD",
		"userList":    n.userCache.Values(),
	}
	n.webSocketRepository.Broadcast(userListPayload)
}
