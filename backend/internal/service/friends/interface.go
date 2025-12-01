package friends

// FriendDTO は handler に返すためのシンプルな DTO です。
// 実際のフィールドは、既存 FriendsHandler が返している JSON に合わせて
// 後続フェーズで調整します（ここでは最小例）。
type FriendDTO struct {
	ID       uint   `json:"ID"`
	Username string `json:"username"`
}

// Service はフレンド関連ユースケースを表すアプリケーションサービスのインターフェースです。
// このフェーズでは「メソッドシグネチャを定義するだけ」で、実装はまだ行いません。
type Service interface {
	// ログインユーザーのフレンド一覧を取得
	ListFriends(userID uint) ([]FriendDTO, error)

	// フレンド申請の作成
	RequestFriend(fromUserID, toUserID uint) error

	// フレンド申請の承認
	ApproveFriend(requestID uint, approverID uint) error
}


