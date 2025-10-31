package friendsHandler

import (
	"fmt"
	"net/http"
	"akichat/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type FriendRequestHandler struct {
	FriendRequestRepo *repository.FriendRequestRepository
}

func NewFriendRequestHandler(friendRequestRepo *repository.FriendRequestRepository) *FriendRequestHandler {
	return &FriendRequestHandler{FriendRequestRepo: friendRequestRepo}
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
	// 型と値のログ出力
	fmt.Printf("userID value: %v, type: %T\n", rawUserID, rawUserID)
	
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

	fmt.Printf("Creating friend request from user %d to user %d\n", userID, req.FriendID)
	if err := h.FriendRequestRepo.CreateFriendRequest(userID, req.FriendID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "リクエストが正常に処理されないかすでにフレンド申請を送っています"})
		return
	}

	if err := h.NotifyFriendRequest(userID,req.FriendID); err != nil {
		fmt.Println("通知エラー",err)
	}else{
		fmt.Println("通知成功")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent successfully"})
}