package UserHandler

import (
	"net/http"
	authsvc "akichat/backend/internal/service/auth"
	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	AuthService authsvc.Service
}

func NewLoginHandler(authService authsvc.Service) *LoginHandler {
	return &LoginHandler{AuthService: authService}
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

	out, err := h.AuthService.Login(authsvc.LoginInput{
		Email: loginData.Email,
		Password: loginData.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// 統一したクッキー名で保存
	c.SetCookie(
		"refreshToken",     // クッキー名（統一）
		out.RefreshToken,   // 保存する値
		60*60*24*7,         // 有効期限（例: 7日間）
		"/",                // パス
		"172.20.10.2",        // ドメイン
		false,              // Secure（HTTPSのみにする）
		true,               // HttpOnly（JSからアクセス不可）
	)

	// パスワードは返却しない。互換性のため username/name/email を付与
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"accessToken": out.AccessToken,
		"userID": out.UserID,
		"username": out.UserName,
		"name": out.UserName,
		"email": out.UserEmail,
	})
}