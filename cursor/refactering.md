了解！フェーズ5は 「router と DI コンテナを整えるフェーズ」 なので、ここまでの分割を前提に、Cursor にそのまま渡せる指示をまとめるね。

# Akichat リファクタリング計画 – フェーズ5 指示（Cursor 用）

## 🧱 前提

このフェーズ5は、以下が完了している前提で進めてください：

- フェーズ0：責務マッピング
- フェーズ1：
  - `internal/service/auth/interface.go`
  - `internal/service/friends/interface.go`
  - `internal/realtime/realtime.go`（`Gateway`）
  - `hub.go` に `var _ realtime.Gateway = (*Hub)(nil)`
- フェーズ2：
  - `internal/service/auth/service.go` で AuthService 実装済み
  - Login / Register handler は `auth.Service` 経由で動いている
- フェーズ3：
  - `internal/service/friends/service.go` で FriendsService 実装済み
  - Friends 系 handler は `friends.Service` 経由で動いている
- フェーズ4：
  - `internal/service/signaling/service.go` で `signaling.Service` 実装済み
  - `Client` 構造体に `Signaling *signaling.Service` が追加され、WebSocket 接続時に注入されている
  - `Client.readPump()` 内の WebRTC シグナリング判定ロジックは `Signaling.Handle(...)` に委譲済み

### 共通制約（フェーズ5も継続）

- HTTP エンドポイント仕様（URL / メソッド / JSON）は変えない。
- WebSocket エンドポイント・メッセージ形式は変えない。
- DB スキーマは変えない。
- 挙動（ログイン、フレンド機能、通話・シグナリングなど）は変えない。

---

## 🎯 フェーズ5のゴール

1. `internal/app/container.go` に「依存関係を組み立てるコンテナ」を追加する。
2. `router.go` から、DB / Repository / Service / Hub の new ロジックをできる限り `Container` に寄せる。
3. handler / WebSocket 用ハンドラは、Container から依存を受け取って動くようにする。
4. `router.go` は「ルーティング定義＋ミドルウェア設定」に集中する形に整理される。

---

## 1. DI コンテナ `internal/app/container.go` の追加

### 1-1. ファイル作成

新規ファイル：

- パス: `backend/internal/app/container.go`

### 1-2. Container 構造体の定義

Akichat で現在扱っている主な依存関係は：

- DB 接続
- 各 Repository
  - `UserRepository`
  - `FriendShipRepository`
  - `FriendRequestRepository`
- 各 Service
  - `auth.Service`
  - `friends.Service`
  - `signaling.Service`
- Realtime（WebSocket）
  - `Hub`（`realtime.Gateway` として扱う）

これらを 1 箇所で組み立てる `Container` を定義します。

```go
package app

import (
    "akichat/internal/db"
    "akichat/internal/handler/webSocket"
    "akichat/internal/realtime"
    "akichat/internal/repository"
    "akichat/internal/service/auth"
    "akichat/internal/service/friends"
    "akichat/internal/service/signaling"
)

type Container struct {
    // インフラ系
    DB  *db.DBType          // 実際の DB 型に合わせて定義
    Hub *webSocket.Hub      // 実装に合わせて型を指定
    RT  realtime.Gateway    // Hub を Gateway として扱う

    // Repository
    UserRepo         repository.UserRepository
    FriendShipRepo   repository.FriendShipRepository
    FriendRequestRepo repository.FriendRequestRepository

    // Service
    AuthService      auth.Service
    FriendsService   friends.Service
    SignalingService *signaling.Service
}


⚠ db.DBType / repository.UserRepository などの型名は、実際のコードに合わせて修正してください。
既存で使っている DB ハンドル（例：*gorm.DB など）や Repo インターフェースの型シグネチャに合わせること。

1-3. NewContainer の実装

NewContainer で依存を組み立てます。
現状 router.go や main.go でやっている初期化処理をすべて集約するイメージです。

func NewContainer() (*Container, error) {
    // 1. DB 初期化（既存 db.go に合わせる）
    database, err := db.NewDB()
    if err != nil {
        return nil, err
    }

    // 2. Hub 初期化
    hub := webSocket.NewHub()
    // Hub.Run() を別 goroutine で起動する必要があれば、
    // Container を返した後に main / server 側で呼ぶか、ここで go hub.Run() してもよい

    // 3. Repository 初期化
    userRepo := repository.NewUserRepository(database)
    friendShipRepo := repository.NewFriendShipRepository(database)
    friendRequestRepo := repository.NewFriendRequestRepository(database)

    // 4. Service 初期化
    authService := auth.NewService(userRepo)
    friendsService := friends.NewService(friendShipRepo, friendRequestRepo, hub)
    signalingService := &signaling.Service{
        RT: hub,
    }

    c := &Container{
        DB:               database,
        Hub:              hub,
        RT:               hub,
        UserRepo:         userRepo,
        FriendShipRepo:   friendShipRepo,
        FriendRequestRepo: friendRequestRepo,
        AuthService:      authService,
        FriendsService:   friendsService,
        SignalingService: signalingService,
    }

    return c, nil
}


⚠ db.NewDB(), webSocket.NewHub(), repository.NewUserRepository() などの実体は
実際の実装に合わせて名称＆引数を修正してください。
すでに GlobalHub が存在する場合、そちらに寄せるか、NewHub を導入して GlobalHub = container.Hub とする構成にするかは、既存設計に合わせて OK です。

2. WebSocket ハンドラを DI フレンドリーにする（必要に応じて）

フェーズ4 では WebSocketHandler が関数形式になっていて、
内部で SignalingService を new していた状態かもしれません。

フェーズ5では、WebSocket ハンドラも Container から依存をもらう構造に寄せます。

2-1. WebSocketHandler 構造体の導入

対象ファイル：

backend/internal/handler/webSocket/websocket_handler.go

例：

package websocket

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/gin-contrib/sessions"

    "akichat/internal/service/signaling"
)

type WebSocketHandler struct {
    Hub       *Hub
    Signaling *signaling.Service
}

func NewWebSocketHandler(hub *Hub, sig *signaling.Service) *WebSocketHandler {
    return &WebSocketHandler{
        Hub:       hub,
        Signaling: sig,
    }
}


※ すでに GlobalHub を使っている場合も、ここでは コンストラクタで受け取る形 を優先します（内部で Global を代入してもOK）。

2-2. ハンドラメソッドとしての WebSocketHandler

これまで func WebSocketHandler(c *gin.Context) だったものを、
構造体メソッドに変えます（ルーティング時は wsHandler.Handle を渡す想定）。

func (h *WebSocketHandler) Handle(c *gin.Context) {
    // 既存の実装をベースにする
    // 1. セッションから user_id を取り出し
    session := sessions.Default(c)
    // ... userID 取得ロジック（既存コードのまま）

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        // 既存のエラーハンドリング
        return
    }

    client := &Client{
        UserID:   userID,
        Conn:     conn,
        Send:     make(chan interface{}, 64),
        Stop:     make(chan struct{}),
        Hub:      h.Hub,
        Signaling: h.Signaling,
    }

    h.Hub.register <- client

    go client.writePump()
    go client.readPump()
}


重要：

ここで SignalingService を new せず、コンテナから渡されたものを使う こと

Hub も GlobalHub を直接参照せず、h.Hub を使うこと

3. router で Container を使って依存を注入する
3-1. router.go で Container を初期化

対象ファイル：

backend/internal/http/router.go

これまで DB / Repo / Service / Hub を直接 new していた処理を、
コンテナ生成に置き換えます。

package http

import (
    "log"

    "github.com/gin-gonic/gin"
    "akichat/internal/app"
    "akichat/internal/handler/userHandler"
    "akichat/internal/handler/friendsHandler"
    "akichat/internal/handler/webSocket"
    // ミドルウェア関連の import もそのまま
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // 1. Container を生成
    container, err := app.NewContainer()
    if err != nil {
        log.Fatalf("failed to init container: %v", err)
    }

    // 2. ここで Hub.Run() を起動する場合もある
    go container.Hub.Run()

    // 3. Handler を DI して生成
    loginHandler := userHandler.NewLoginHandler(container.AuthService)
    registerHandler := userHandler.NewRegisterHandler(container.AuthService)

    friendsHandler := friendsHandler.NewFriendsHandler(container.FriendsService)
    requestHandler := friendsHandler.NewFriendRequestHandler(container.FriendsService)
    approveHandler := friendsHandler.NewApproveRequestHandler(container.FriendsService)

    wsHandler := webSocket.NewWebSocketHandler(container.Hub, container.SignalingService)

    // 4. ミドルウェア設定（CORS、セッション、JWT など）
    //    → 既存の設定をそのまま残しつつ、必要なら container の情報を利用

    // 5. ルーティング定義
    api := r.Group("/api")
    {
        api.POST("/login", loginHandler.Handle)
        api.POST("/register", registerHandler.Handle)

        api.GET("/friends", friendsHandler.Handle)
        api.POST("/friends/request", requestHandler.Handle)
        api.POST("/friends/approve", approveHandler.Handle)

        api.GET("/websocket", wsHandler.Handle) // 実際のパスに合わせる
    }

    return r
}


⚠ ルートパス・グルーピング・HTTP メソッド名は、実際のコードに合わせて調整してください。
重要なのは「router の中で new していた Repository / Service / Hub の組み立てを、Container に寄せる」ことです。

3-2. router から「組み立てロジック」を減らす

SetupRouter 内に残してよいもの：

Gin エンジンの生成

ミドルウェア設定（CORS / ログ / リカバリ / セッション）

エンドポイントと handler の紐付け

SetupRouter から減らしたいもの：

DB 接続構築（sql.Open, gorm.Open, db.NewDB 等）

Repository の New...Repository(...)

Service の New...Service(...)

Hub := NewHub() 的な生成

すべてコンテナ側に寄せることで、router はフレームワークに近い層 になります。

4. main.go / server 起動部分の調整（必要なら）

もし cmd/server/main.go で SetupRouter() を呼んでいる構成なら、
そこは基本的に今まで通りで OK です。

例：

func main() {
    r := http.SetupRouter()
    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}


Container の生成や Hub.Run() は、SetupRouter 内で完結している設計にしてよいですし、
逆に main 側で：

container, _ := app.NewContainer()

hub := container.Hub; go hub.Run()

router := http.SetupRouter(container)

のように「Container を main で作って router に渡す」スタイルに変えるのもアリです。

このフェーズ5では、どちらかに統一されていれば OK とします。

✅ フェーズ5の終了条件

go build ./... が成功すること。

backend/internal/app/container.go が存在し、主要な依存（DB / Repo / Service / Hub / signaling）がそこに集約されていること。

router.go から：

DB 初期化、Repository の new、Service の new、Hub の new といった処理がほぼ消えていること。

代わりに Container を生成し、そこから Handler / WebSocketHandler を new していること。

WebSocketHandler / FriendsHandler / LoginHandler / RegisterHandler などが：

直接 Repository や Hub を new せず、コンストラクタ引数で依存（Service / Hub / Signaling） を受け取っていること。

アプリを起動して、以下がフェーズ5前と同じように動くこと：

ログイン / 登録

フレンド一覧 / 申請 / 承認 ＋ 通知

WebSocket 接続 ＋ WebRTC シグナリング（通話開始〜ICE 交換）

これでフェーズ1〜5を通じて、

handler = 入出力（HTTP/WebSocket）担当

service = ユースケース担当

repository = DB 担当

realtime = 通信インフラ担当

container = 依存組み立て担当

というキレイな責務分離に近づきます。