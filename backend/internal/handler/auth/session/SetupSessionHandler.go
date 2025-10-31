package session

import (
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetupSessionHandler(c *gin.Context) {
	session := sessions.Default(c)

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// セッションに許可フラグとユーザーIDを保存
	session.Set("websocket_allowed", true)
	session.Set("user_id", userID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session save error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "WebSocket allowed"})
}