package friendsHandler

import (
	"net/http"
	svc "akichat/backend/internal/service/friends"
	"github.com/gin-gonic/gin"
)

type FriendShipHandler struct {
	FriendsService svc.Service
}

func NewFriendShipHandler(service svc.Service) *FriendShipHandler {
	return &FriendShipHandler{FriendsService: service}
}

func (h *FriendShipHandler) GetFriendsHandler(c *gin.Context) {
	// ユーザーIDをコンテキストから取得（ミドルウェアで設定されていることを前提）
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	friends, err := h.FriendsService.ListFriends(userIDInterface.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve friends"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"friends": friends})
}
