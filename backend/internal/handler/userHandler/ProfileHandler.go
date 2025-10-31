package UserHandler

import (
	"net/http"
	"akichat/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	//Repoという変数にはuserrepository型のポインタを格納するという定義
	Repo *repository.UserRepository
}

func NewProfileHandler(repo *repository.UserRepository) *ProfileHandler {
	return &ProfileHandler{Repo: repo}
}

//メソッドレシーバーは関数に引数として渡すこともできる
func (h *ProfileHandler) ProfileHandler(c *gin.Context) {
	// 認証されたユーザーの情報を取得
	userEmail, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// ユーザー情報をデータベースから取得
	//h.Repoはuserrepository型のポインタ
	user, err := h.Repo.GetUserByEmail(userEmail.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email": user.Email,
		"name":  user.Username,
		// 他の必要なユーザー情報を追加
	})
}