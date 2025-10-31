package friendsHandler

import (
	"fmt"
	"net/http"
	"akichat/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type SearchNotFriendHandler struct {
	UserRepo       *repository.UserRepository
}

func NewSearchNotFriendHandler(userRepo *repository.UserRepository) *SearchNotFriendHandler {
	return &SearchNotFriendHandler{UserRepo: userRepo}
}
func (h *SearchNotFriendHandler) SearchNotFriendHandler(c *gin.Context) {
	fmt.Printf("start to search")
	var req struct {
		Keyword string `json:"keyword"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		fmt.Printf("invalid keyword")
		return
	}

	if req.Keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keyword is required"})
		return
	}

	user, err := h.UserRepo.GetUserByUserName(req.Keyword)
	if err != nil {
		fmt.Printf("user not found")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user"})
		return
	}
	fmt.Printf(user.Username);
	c.JSON(http.StatusOK, gin.H{"user": user})
}
