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

    client := &Client{
		UserID:   userID,
		Conn: conn,
		Send: make(chan interface{}, 64), // バッファ推奨
		Stop: make(chan struct{}),
		Hub:  GlobalHub,
	}

	GlobalHub.register <- client

	go client.writePump()
	go client.readPump()
}
