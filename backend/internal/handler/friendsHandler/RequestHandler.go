package friendsHandler

import (
	"net/http"
	svc "akichat/backend/internal/service/friends"
	"github.com/gin-gonic/gin"
)

type FriendRequestHandler struct {
	FriendsService svc.Service
}

func NewFriendRequestHandler(service svc.Service) *FriendRequestHandler {
	return &FriendRequestHandler{FriendsService: service}
}

func (h *FriendRequestHandler) FriendRequestHandler(c *gin.Context) {
	var req struct {
		FriendID uint `json:"to_user_id"`
	}
	rawUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := rawUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.FriendID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "FriendID is required"})
		return
	}

	if err := h.FriendsService.RequestFriend(userID, req.FriendID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "リクエストが正常に処理されないかすでにフレンド申請を送っています"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent successfully"})
}