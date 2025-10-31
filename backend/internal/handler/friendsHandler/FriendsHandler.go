package friendsHandler

import (
	"net/http"
	"akichat/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type FriendShipHandler struct {
	FriendShipRepo *repository.FriendShipRepository
}

func NewFriendShipHandler(friendShipRepo *repository.FriendShipRepository) *FriendShipHandler {
	return &FriendShipHandler{FriendShipRepo: friendShipRepo}
}

func (h *FriendShipHandler) GetFriendsHandler(c *gin.Context) {
	// ユーザーIDをコンテキストから取得（ミドルウェアで設定されていることを前提）
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	friends, err := h.FriendShipRepo.GetFriendsByUserID(userIDInterface.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve friends"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"friends": friends})
}
