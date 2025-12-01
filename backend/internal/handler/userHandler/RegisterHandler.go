package UserHandler

import (
	"net/http"
	authsvc "akichat/backend/internal/service/auth"
	"akichat/backend/internal/model"
	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	AuthService authsvc.Service
}

func NewRegisterHandler(authService authsvc.Service) *RegisterHandler {
	return &RegisterHandler{AuthService: authService}
}

func (h *RegisterHandler) RegisterHandler(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	out, err := h.AuthService.Register(authsvc.RegisterInput{
		Name: user.Username,
		Email: user.Email,
		Password: user.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.SetCookie(
		"refreshToken",       // クッキー名
		out.RefreshToken,     // 保存する値
		60*60*24*7,           // 有効期限（例: 7日間）
		"/",                  // パス
		"172.20.10.2",          // ドメイン（本番環境では自ドメインを指定）
		false,                 // Secure（HTTPSのみにする）
		true,                 // HttpOnly（JSからアクセス不可）
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"accessToken": out.AccessToken, 
		"id": out.UserID,
		"userID": out.UserID, 
		"name": out.UserName,
		"username": out.UserName, 
		"email": out.UserEmail,
	})
}

