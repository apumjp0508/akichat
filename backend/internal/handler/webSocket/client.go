package websocket

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	signaling "akichat/backend/internal/service/signaling"
)

const (
	writeWait  = 10 * time.Second
	pingPeriod = 30 * time.Second
	readWait   = 60 * time.Second         
)

type Client struct {
	UserID uint
	Conn   *websocket.Conn
	Send   chan interface{}
	Stop chan struct{}
	Hub  *Hub
	Signaling *signaling.Service
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub側でcloseされた場合
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// JSONで送信（直列化されているので安全）
			if err := c.Conn.WriteJSON(msg); err != nil {
				fmt.Println("writePump WriteJSON err:", err)
				return
			}

		case <-ticker.C:
			// 定期Ping（PongはRead側 or Gorilla内部で処理）
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("writePump ping err:", err)
				return
			}

		case <-c.Stop:
			return
		}
	}
}

// クライアント→サーバ：メッセージ受信＋切断検知
func (c *Client) readPump() {
	defer func() {
		// 終了時のクリーンアップをHubへ依頼
		c.Hub.unregister <- c
	}()

	// あなたの setupingpong.go の関数を呼ぶ（ReadDeadline + PongHandler）
	SetupPingPong(c.Conn, int(readWait/time.Second))

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			// Close/Ping/Pongタイムアウト等はここに入ってくる
			// 例: websocket: close 1001 (going away)
			return
		}

		// 1) シグナリング処理をサービスに委譲
		if c.Signaling != nil {
			if err := c.Signaling.Handle(c.UserID, msg); err != nil {
				// 既存の webrtc_error レスポンスに合わせて通知
				if de, ok := err.(*signaling.DeliveryError); ok {
					select {
					case c.Send <- map[string]interface{}{
						"type":         "webrtc_error",
						"reason":       "user_offline",
						"to":           de.To,
						"originalType": de.OriginalType,
					}:
					case <-c.Stop:
						return
					}
				}
			}
		}

		// 任意：ACKなどのアプリメッセージをここで処理
		if string(msg) == "notification_ack" {
			select {
			case c.Send <- map[string]string{
				"type":    "ack_confirmed",
				"message": "notification received successfully",
			}:
			case <-c.Stop:
				return
			}
		}
	}
}
