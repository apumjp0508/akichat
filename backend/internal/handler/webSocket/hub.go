package websocket

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// ユーザーIDごとに接続を保持
	Clients map[uint]*websocket.Conn
	mu      sync.Mutex
}

//ここでポインタ型のインスタンスを生成している
var GlobalHub = &Hub{
	//：は構造体 Hub のフィールド Clients 
	// に、make(map[uint]*websocket.Conn) という値を代入するという意味
	//makeはからの連想配列を作成する
	Clients: make(map[uint]*websocket.Conn),
	// 	Gorilla WebSocketの接続オブジェクト
	// ユーザーがWebSocketで接続したときの「接続そのもの」を表す構造体へのポインタ
	// 	｛例｝Clients := map[uint]*websocket.Conn{
	//     1001: connA,
	//     2005: connB,
	// }

}

// 接続を登録
func (h *Hub) Register(userID uint, conn *websocket.Conn) {
	//複数のuserが同時に通知を送ろうとしてmapに同時アクセスしてしまうとエラーが起きる
	//h.mu/Lock()でほかからのアクセスをブロックすることができる
	h.mu.Lock()
	//deferを使うと処理中に panic や return が起きても、確実に Unlock() を呼べる
	//関数の処理が全て終わったあとに呼び出す
	defer h.mu.Unlock()
	h.Clients[userID] = conn
	fmt.Printf("✅ User %d registered WebSocket connection\n", userID)
}

// 接続を削除
func (h *Hub) Unregister(userID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conn, ok := h.Clients[userID]; ok {
		conn.Close()
		delete(h.Clients, userID)
		fmt.Printf("❌ User %d connection closed\n", userID)
	}
}

// 接続中の全ユーザーIDを返す
func (h *Hub) GetConnectedUsers() []uint {
	h.mu.Lock()
	defer h.mu.Unlock()

	userIDs := make([]uint, 0, len(h.Clients))
	for id := range h.Clients {
		userIDs = append(userIDs, id)
	}
	return userIDs
}

// 特定ユーザーに通知を送る
func (h *Hub) NotifyUser(requestUserID uint, userID uint, message string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	//ここでclients連想配列にアクセスuseridに相当するwebsocket通信を取得する
	conn, ok := h.Clients[userID]
	if !ok {
		return fmt.Errorf("user %d is not connected", userID)
	}

	//ここで実際に通信を行いjson形式で相手にメッセージを送信する
	err := conn.WriteJSON(map[string]interface{}{
		"type":    "friend_request",
		"message": message,
		"requestUserID": requestUserID,
	})
	if err != nil {
		conn.Close()
		delete(h.Clients, userID)
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}
