package UserHandler

import (
	"fmt"
	"net/http"
	"akichat/backend/internal/repository"
	"akichat/backend/internal/handler/auth/token/JWTToken"
	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	Repo *repository.UserRepository
}

func NewLoginHandler(repo *repository.UserRepository) *LoginHandler {
	return &LoginHandler{Repo: repo}
}

func (h *LoginHandler) LoginHandler(c *gin.Context) {
	var loginData struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}


    // ここでユーザー認証を行う（例: データベースと照合）
    user, err := h.Repo.GetUserByEmail(loginData.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}


	accessToken, refreshToken, err := JWTHandler.GenerateTokens(user.ID, user.Email)
	if err != nil {
		fmt.Println("Error generating token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}else {	
		fmt.Println("Generated Token:", accessToken)
	}

	// 統一したクッキー名で保存
	c.SetCookie(
		"refreshToken",     // クッキー名（統一）
		refreshToken,       // 保存する値
		60*60*24*7,         // 有効期限（例: 7日間）
		"/",                // パス
		"localhost",        // ドメイン
		false,              // Secure（HTTPSのみにする）
		true,               // HttpOnly（JSからアクセス不可）
	)


	//未使用の変数を一時的に対比するために使用している
    _ = user
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "accessToken": accessToken, "userID":user.ID, "username":user.Username, "password":user.Password})
}