package websocket

import (
	"fmt"
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

func StartPingLoop(conn *websocket.Conn, intervalSec int, stopCh <-chan struct{}) {
	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer ticker.Stop()

	//for単体で使用すると無限ループになる
	for {
		//selectはcaseを一つだけ実行できる
		select {
		case <-ticker.C:
			// 一定間隔でPing送信
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("Ping送信失敗:", err)
				conn.Close()
				return
			}
		case <-stopCh:
			// 明示的に停止指示が来たとき
			return
		}
	}
}
