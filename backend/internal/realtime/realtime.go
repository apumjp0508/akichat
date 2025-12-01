package realtime

// Gateway は、アプリケーション層から見た
// 「リアルタイム通信（WebSocket）の窓口」を表すポートです。
// 現時点では、既存の Hub 実装に合わせて定義します。
type Gateway interface {
	// 接続中のユーザーID一覧を取得する
	GetConnectedUsers() []uint

	// 指定したユーザーにフレンドリクエストなどの通知メッセージを送る
	// （既存 Hub.NotifyUser のシグネチャに合わせる）
	NotifyUser(requestUserID, userID uint, message string) error

	// 任意のペイロードを指定ユーザーに送信する汎用メソッド
	// （既存 Hub.SendTo のシグネチャに合わせる）
	SendTo(userID uint, payload interface{}) error
}



