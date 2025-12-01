package auth

// LoginInput はログイン時の入力を表す DTO です。
type LoginInput struct {
	Email    string
	Password string
}

// LoginOutput はログイン成功時に返す情報です。
// 実際のレスポンス仕様に合わせてフィールドを調整してください。
type LoginOutput struct {
	AccessToken  string
	RefreshToken string

	UserID    uint
	UserName  string
	UserEmail string
}

// RegisterInput はユーザー登録時の入力を表す DTO です。
type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

// RegisterOutput は登録成功時に返す情報です。
type RegisterOutput struct {
	AccessToken  string
	RefreshToken string

	UserID    uint
	UserName  string
	UserEmail string
}

// Service は認証・登録に関するユースケースを表すインターフェースです。
// このフェーズでは、まだ実装は追加せずシグネチャのみ定義します。
type Service interface {
	Login(in LoginInput) (LoginOutput, error)
	Register(in RegisterInput) (RegisterOutput, error)
}


