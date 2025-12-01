package friendsHandler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func (h *FriendShipHandler) ApproveRequestHandler(c *gin.Context) {
	var req struct {
		RequestUserID uint `json:"requestUserID"`
		UserID        uint `json:"userID"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.FriendsService.ApproveFriend(req.RequestUserID, req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "フレンド追加に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Add friend successfully"})
}
