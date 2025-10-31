package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gin-contrib/sessions"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(c *gin.Context) {
	session := sessions.Default(c)
	userIDVal := session.Get("user_id")
	userID, ok := userIDVal.(uint)
	if err != nil {
		fmt.Println("Invalid userid")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	//Upgrader()はHTTP通信をwebsocket通信に変換するメソッド
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("failed to upgrade connection")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	GlobalHub.Register(uint(userID), conn)

	// 受信ループ（クライアント→サーバーのACK用）
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			GlobalHub.Unregister(uint(userID))
			break
		}
		// クライアントから通知確認メッセージを受け取る
		if string(msg) == "notification_ack" {
			conn.WriteJSON(map[string]string{
				"type":    "ack_confirmed",
				"message": "notification received successfully",
			})
		}
	}
}
