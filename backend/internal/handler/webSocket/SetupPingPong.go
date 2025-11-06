package websocket

import (
	"time"
	"github.com/gorilla/websocket"
)

func SetupPingPong(conn *websocket.Conn, timeoutSec int) {
	timeout := time.Duration(timeoutSec) * time.Second

	conn.SetReadDeadline(time.Now().Add(timeout))

	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(timeout))
		return nil
	})
}
