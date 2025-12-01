# 責務マッピング（フェーズ0）

このドキュメントは現状の挙動を一切変更せず、主要モジュールの「理想の責務」「実際にやっていること」「責務外と思われる処理」を整理したものです。改善は次フェーズ以降のため、本書では事実ベースの観察を優先します。

---

## internal/handler/userHandler/LoginHandler.go

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） | ログインリクエストの受領・バリデーション、アプリ層への委譲、トークン等の結果をHTTPレスポンスに変換 |
| 実際にやっていること | JSONバインド、リポジトリでユーザー検索、JWT発行、refreshTokenクッキー設定（固定ドメイン）、JSONレスポンス返却 |
| 明らかに責務外と思われる処理 | JWT発行の詳細、クッキー属性/ドメインの直接指定、パスワードをレスポンスへ含める挙動（情報露出） |

---

## internal/handler/userHandler/RegisterHandler.go

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） | 登録リクエストの受領・バリデーション、アプリ層へユーザー作成を委譲、結果をHTTPレスポンスに変換 |
| 実際にやっていること | JSONバインド、リポジトリでユーザー作成、JWT発行、refreshTokenクッキー設定（固定ドメイン）、JSONレスポンス返却 |
| 明らかに責務外と思われる処理 | JWT発行の詳細、クッキー属性/ドメインの直接指定、パスワードをレスポンスへ含める挙動（情報露出） |

---

## internal/handler/userHandler/ProfileHandler.go（GetMeHandler）

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） | 認証済みユーザーのプロフィール取得の入口、アプリ層に照会しHTTPレスポンスを構築 |
| 実際にやっていること | ContextからuserID取得、リポジトリでユーザー検索、id/email/name をJSONで返却 |
| 明らかに責務外と思われる処理 | 取得ロジックの直接呼び出し（アプリ層が薄い/存在しない） |

---

## internal/handler/friendsHandler/FriendsHandler.go（FriendShipHandler）

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） | フレンド一覧取得リクエストの入口、アプリ層経由で一覧を取得しHTTP化 |
| 実際にやっていること | ContextのuserID取得、FriendShipRepoで一覧を取得、JSON返却 |
| 明らかに責務外と思われる処理 | 取得ロジックの直接呼び出し（アプリ層が薄い/存在しない） |

---

## internal/handler/friendsHandler/RequestHandler.go（FriendRequestHandler）

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） | フレンド申請の受付、アプリ層へ作成依頼、結果をHTTP化 |
| 実際にやっていること | ContextのuserID取得、JSONバインド、FriendRequestRepoで申請作成、WebSocket通知呼び出し、JSON返却 |
| 明らかに責務外と思われる処理 | 通知送信（WebSocket）まで直接呼び出し（アプリ/ドメイン層へ委譲したい） |

---

## internal/handler/friendsHandler/ApproveRequestHandler.go

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） | フレンド承認の受付、アプリ層で承認処理し結果をHTTP化 |
| 実際にやっていること | JSONバインド、FriendShipRepo.AddFriend呼び出し、結果をJSON返却 |
| 明らかに責務外と思われる処理 | リポジトリ直接呼び出し（アプリ層が薄い/存在しない） |

---

## internal/handler/friendsHandler/NotifyFriendRequest.go

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） |（アプリ層の）通知ユースケースから呼ばれるインフラ通知ゲートウェイ |
| 実際にやっていること | Hubへ直接アクセスして対象ユーザーへ通知メッセージを送信 |
| 明らかに責務外と思われる処理 | ハンドラ層からHubへ直接依存（ポート/アダプタの分離がない） |

---

## internal/handler/webSocket/client.go

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） | 単一WebSocketクライアントの接続・送受信管理、KeepAlive、上位層へのイベント伝達 |
| 実際にやっていること | write/readポンプ、Ping送信/ReadDeadline、WebRTCシグナリングJSONの受信とHubへの配送、オフライン時のエラー返信、ACK応答 |
| 明らかに責務外と思われる処理 | シグナリングのメッセージ種別判定まで保持（分離自体は許容範囲だが、ユースケース層が薄い） |

---

## internal/handler/webSocket/hub.go

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） | 接続クライアントのレジストリ、メッセージ配送、接続管理 |
| 実際にやっていること | register/unregister、クライアント再接続時の古い接続クローズ、特定ユーザーへの配送、接続中ユーザー一覧の取得、通知送信 |
| 明らかに責務外と思われる処理 | なし（Hubの責務に概ね収まっている） |

---

## internal/handler/webSocket/websocket_handler.go

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） | HTTP→WebSocketアップグレードの入口、認証/セッションの検証、接続のHub登録 |
| 実際にやっていること | セッションからuser_id抽出、GorillaでUpgrade、Client構築とHub登録、read/writeポンプ起動 |
| 明らかに責務外と思われる処理 | なし（HTTP処理の範囲内） |

---

## internal/handler/webSocket/getConnectedUsersHandler.go

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） | 現在接続中のユーザー一覧取得のHTTPエンドポイント |
| 実際にやっていること | Hubから接続中ユーザーID一覧を取得しJSON返却 |
| 明らかに責務外と思われる処理 | なし |

---

## internal/handler/webSocket/SetupPingPong.go

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） | WebSocketのReadDeadline/PongHandlerを安全に設定するユーティリティ |
| 実際にやっていること | 指定タイムアウトでReadDeadline設定、Pongで延長 |
| 明らかに責務外と思われる処理 | なし |

---

## internal/http/router.go

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） | ルータ初期化、ミドルウェア適用、DIの組立（Repo/Handler）、エンドポイント定義 |
| 実際にやっていること | DB初期化、Config読込、CookieセッションStore設定、CORS設定、Repo/Handler生成、API/WSルーティング定義 |
| 明らかに責務外と思われる処理 | 単純なDI以上のビルド（複雑化の予兆はあるが現状は許容範囲） |

---

## internal/service/user-service.go

| 項目 | 内容 |
|------|------|
| 本来の責務（理想） | ユーザー関連のユースケース（アプリケーションサービス）を実装 |
| 実際にやっていること | 現状ファイルのみ（未実装） |
| 明らかに責務外と思われる処理 | なし |

---

## 依存関係と外部インターフェース（概観）

- 認証系: JWT発行（handlerから直接JWTモジュールを呼び出し）、HTTPクッキーにrefreshToken設定（ドメイン固定/セキュア属性は環境依存のため要見直し余地）  
- セッション: Cookieセッション（`/api/session/websocket` のアップグレードで使用）  
- WebSocket: Gorillaベース。`Hub` がクライアントをID単位で管理し、`Client` がread/writeポンプ/KeepAliveを担当。シグナリング（`webrtc_offer/answer/ice`）はJSONで受け、宛先ユーザーに中継。  
- DBアクセス: Repository直呼び出し（handler→repo が多く、アプリ層は薄い）  
- CORS/セッション設定: `config.Load()` の値で制御。  

---

## 次フェーズ以降の改善メモ（参考）

- ハンドラからJWT発行・通知（WS）などのロジックを直接呼ばず、アプリ層（service）へ委譲する（ポート/アダプタ分離）。  
- レスポンスの機微情報（例: パスワード）の排除。  
- クッキー属性（`Secure`/`SameSite`/`Domain`）を環境に応じて統一管理。  
- WebSocketシグナリングは専用のアプリ層/ユースケースを定義し、Hubは輸送レイヤに限定。  


