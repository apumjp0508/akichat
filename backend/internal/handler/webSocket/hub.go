package websocket

import (
	"fmt"
)

type Hub struct {
	// ユーザーIDごとに接続を保持
	clients map[uint]*Client
	register   chan *Client
	unregister chan *Client
}

var GlobalHub = NewHub()

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			// 既に同じIDが居たら古い接続を落とす（再接続対策）
			if old, ok := h.clients[c.UserID]; ok {
				close(old.Stop)
				close(old.Send)
				old.Conn.Close()
			}
			h.clients[c.UserID] = c

		case c := <-h.unregister:
			if cur, ok := h.clients[c.UserID]; ok && cur == c {
				delete(h.clients, c.UserID)
				// SendをcloseするとwritePumpが終わる
				close(c.Send)
				close(c.Stop)
				c.Conn.Close()
			}
		}
	}
}

// 接続中の全ユーザーIDを返す
func (h *Hub) GetConnectedUsers() []uint {
	ids := make([]uint, 0, len(h.clients))
	//rangeはchに値が送信されるたびに一回ループが回る
	for id := range h.clients {
		ids = append(ids, id)
	}
	return ids
}

func (h *Hub) NotifyUser(requestUserID, userID uint, message string) error {
	c, ok := h.clients[userID]
	if !ok {
		return fmt.Errorf("user %d is not connected", userID)
	}
	select {
	case c.Send <- map[string]interface{}{
		"type":          "friend_request",
		"message":       message,
		"requestUserID": requestUserID,
	}:
		return nil
	case <-c.Stop:
		return fmt.Errorf("user %d already closed", userID)
	}
}

// 任意のペイロードを宛先ユーザーに送信する汎用メソッド
//ペイロード（payload）＝実際に送りたい内容（チャット本文や通知内容など）
// writePump が c.Send を受け取り実送信するため、ここではチャネルに投入するだけでよい
func (h *Hub) SendTo(userID uint, payload interface{}) error {
	c, ok := h.clients[userID]
	if !ok {
		return fmt.Errorf("user %d is not connected", userID)
	}
	select {
	case c.Send <- payload:
		return nil
	case <-c.Stop:
		return fmt.Errorf("user %d already closed", userID)
	}
}