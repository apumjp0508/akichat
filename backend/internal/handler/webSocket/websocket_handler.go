package websocket

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gin-contrib/sessions"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(c *gin.Context) {
    fmt.Println("websocket通信を開始")

    session := sessions.Default(c)
    var userID uint
    switch v := session.Get("user_id").(type) {
    case int:
        userID = uint(v)
    case int64:
        userID = uint(v)
    case float64:
        userID = uint(v)
    case uint:
        userID = v
    default:
        fmt.Println("Invalid user_id type:", v)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
        return
    }

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        fmt.Println("failed to upgrade connection:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
        return
    }

    GlobalHub.Register(userID, conn)

    // Ping/Pong設定と開始
    SetupPingPong(conn, 60)
    stopCh := make(chan struct{})
    go StartPingLoop(conn, 30, stopCh)

    // 受信ループ
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Printf("❌ connection closed for user %d: %v\n", userID, err)
            close(stopCh)
            GlobalHub.Unregister(userID)
            break
        }

        if string(msg) == "notification_ack" {
            conn.WriteJSON(map[string]string{
                "type":    "ack_confirmed",
                "message": "notification received successfully",
            })
        }
    }
}
