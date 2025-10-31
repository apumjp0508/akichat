package websocket

import(
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetConnectedUsersHandler(c *gin.Context) {
	userIDs := GlobalHub.GetConnectedUsers()
	c.JSON(http.StatusOK, gin.H{
		"connected_users": userIDs,
	})
}
