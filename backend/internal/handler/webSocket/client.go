package websocket

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
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

		// WebRTC シグナリングのJSONメッセージ受付
		type inboundSignal struct {
			Type      string          `json:"type"`
			To        uint            `json:"to"`
			SDP       json.RawMessage `json:"sdp,omitempty"`
			Candidate json.RawMessage `json:"candidate,omitempty"`
		}
		var m inboundSignal
		//json.UnmarshalはJSONをパースして構造体に格納する関数
		if err := json.Unmarshal(msg, &m); err == nil && m.Type != "" {
			switch m.Type {
			case "webrtc_offer", "webrtc_answer":
				payload := map[string]interface{}{
					"type": m.Type,
					"from": c.UserID,
					"sdp":  m.SDP,
				}
				fmt.Println("webrtc_offer or webrtc_answer:", payload)
				if err := c.Hub.SendTo(m.To, payload); err != nil {
					select {
					case c.Send <- map[string]interface{}{
						"type":         "webrtc_error",
						"reason":       "user_offline",
						"to":           m.To,
						"originalType": m.Type,
					}:
					case <-c.Stop:
						return
					}
				}
				continue
			case "webrtc_ice":
				payload := map[string]interface{}{
					"type":      m.Type,
					"from":      c.UserID,
					"candidate": m.Candidate,
				}
				fmt.Println("webrtc_ice:", payload)
				if err := c.Hub.SendTo(m.To, payload); err != nil {
					select {
					case c.Send <- map[string]interface{}{
						"type":         "webrtc_error",
						"reason":       "user_offline",
						"to":           m.To,
						"originalType": m.Type,
					}:
					case <-c.Stop:
						return
					}
				}
				continue
			default:
				// 未知タイプは無視
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
