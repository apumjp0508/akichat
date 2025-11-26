## 概要
slackのハドルのようなビデオ通話ができるアプリを作成しました。

## 技術スタック
### フロントエンド
| 技術               | バージョン  | 用途                      |
| ---------------- | ------ | ----------------------- |
| **Next.js**      | 15.5.4 | Reactベースのフルスタックフレームワーク  |
| **React**        | 19.1.1 | UI構築ライブラリ               |
| **TypeScript**   | 5.9.2  | 型安全なJavaScript開発        |
| **Tailwind CSS** | 3.4.18 | ユーティリティファーストなCSSフレームワーク |
| **Zustand**      | 5.0.8  | 軽量な状態管理ライブラリ            |
| **ESLint**       | 9.36.0 | 静的解析・コード品質維持ツール         |

###  バックエンド
| 技術                    | バージョン  | 用途                            |
| --------------------- | ------ | ----------------------------- |
| **Go**                | 1.24.0 | バックエンドプログラミング言語               |
| **Gin**               | 1.11.0 | 軽量なWebフレームワーク                 |
| **GORM**              | 1.31.0 | ORM（Object Relational Mapper） |
| **gorilla/websocket** | 1.5.3  | 双方向通信（WebSocket）              |
| **golang-jwt/jwt**    | v5.3.0 | JWT認証ライブラリ                    |

### インフラストラクチャ
| 技術                          | バージョン | 用途              |
| --------------------------- | ----- | --------------- |
| **Docker / Docker Compose** | -     | 開発・本番環境のコンテナ化   |
| **MySQL**                   | 8.x   | 永続データベース        |
| **Redis**                   | 7.x   | キャッシュおよびセッション管理 |

## 音声、映像の通信フロー
<img width="425" height="378" alt="スクリーンショット 2025-11-26 200543" src="https://github.com/user-attachments/assets/ceed16fe-a516-4b86-92f4-019a6536f07a" />

##  WebRTC 接続の流れ

ブラウザ間で音声・映像を直接通信するまでの基本的な WebRTC の処理手順です。  
シグナリングサーバーを介して接続情報を交換し、最終的に P2P 通信が確立されます。

---

### 通信プロセス

**1. A が STUN サーバーを通して自身の接続候補（IP/ポート）情報を取得する**  
→ `RTCPeerConnection` の `icecandidate` イベントで ICE 情報を得る。

---

**2. A がサーバーに offer（SDP）を送信**  
→ サーバーはこれを **シグナリングサーバー（WebSocket 等）** 経由で転送。

---

**3. サーバー経由で B に offer が送られる**

---

**4. B が answer（SDP）を作成しサーバーへ返す**

---

**5. サーバーが A に answer を送信**

---

**6. 双方が ICE candidate（接続候補）を交換（サーバー経由）**

---

**7. PeerConnection が確立**  
→ A ↔ B 間で直接通信できるようになる。

---

**8. A ↔ B 間で音声・映像を直接送受信**  
→ 実際のメディアデータは **P2P（Peer to Peer）通信** でやり取りされ、サーバーを経由しない。

---
# Realtime 層の抽象化リファクタ案（WebSocket / WebRTC 対応）

## ゴール

`akichat` における **WebSocket / WebRTC などインフラ依存コード**を  
`interface` でカプセル化し、次の状態を目指す：

- アプリケーションのユースケース（Chat / Call など）は  
  **「どう送るか」ではなく「誰に何を送りたいか」だけを書く**
- 通信方式（WebSocket / WebRTC DataChannel / それ以外）が変わっても  
  **ユースケース層の変更を最小限に抑えられる**
- テスト時にモック実装へ差し替えやすくする

---

## 目的・メリット

- **疎結合**  
  - アプリ側は `WebSocket` や `PeerConnection` の具体実装を知らずに済む  
- **保守性向上**  
  - 通信方式やインフラ構成の変更（例：WebSocket → WebRTC、別リージョン）に強くなる
- **テスト容易性**  
  - `interface` のモックを使うことで、ユースケースのユニットテストがしやすくなる
- **拡張性**  
  - 将来「WebRTC の DataChannel でチャット」「別プロトコル対応」などがやりやすい

---

## 全体像

### Before（現状イメージ）

- `ChatService` / `Hub` / Handler などが **直接 `*websocket.Conn` や PeerConnection に依存**

```go
// （イメージ）アプリ層が直接 WebSocket に依存している例
func (s *ChatService) SendMessage(userID uint, text string) error {
    // WebSocket のコネクション構造体に直接アクセス
    conn := s.hub.GetConn(userID)
    return conn.WriteJSON(map[string]any{
        "type": "chat",
        "text": text,
    })
}
