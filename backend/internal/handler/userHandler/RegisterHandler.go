package UserHandler

import (
	"fmt"
	"net/http"
	"akichat/backend/internal/repository"
	"akichat/backend/internal/model"
	"akichat/backend/internal/handler/auth/token/JWTToken"
	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	Repo *repository.UserRepository
}

func NewRegisterHandler(repo *repository.UserRepository) *RegisterHandler {
	return &RegisterHandler{Repo: repo}
}

func (h *RegisterHandler) RegisterHandler(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.Repo.CreateUser(&user); err != nil {
		fmt.Println("Error creating user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	accessToken, refreshToken, err := JWTHandler.GenerateTokens(user.ID,user.Email)
	if err != nil {
		fmt.Println("Error generating token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}else {	
		fmt.Println("Generated Token:", accessToken)
	}

	c.SetCookie(
		"refreshToken",       // クッキー名
		refreshToken,         // 保存する値
		60*60*24*7,           // 有効期限（例: 7日間）
		"/",                  // パス
		"172.20.10.2",          // ドメイン（本番環境では自ドメインを指定）
		false,                 // Secure（HTTPSのみにする）
		true,                 // HttpOnly（JSからアクセス不可）
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"accessToken": accessToken, 
		"id": user.ID, 
		"name": user.Username, 
		"email": user.Email, 
		"password": user.Password,
	})
}

