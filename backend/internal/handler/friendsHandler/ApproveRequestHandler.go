package friendsHandler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func (h *FriendShipHandler) ApproveRequest(c *gin.Context) {
	var req struct {
		RequestUserId uint `json:"requestUserId"`
		UserId        uint `json:"userID"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.FriendShipRepo.AddFriend(req.RequestUserId, req.UserId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "フレンド追加に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Add friend successfully"})
}
