## 概要
webrtc通信を実現したいと考えている。

## すでに実装済み
1) Hub に汎用送信を追加 notifiuser()は変更せず新規でメソッド追加する
2) Client.readPump をJSON対応へ拡張
受信JSONを {type,to,sdp,candidate} として Unmarshal
switch type:
webrtc_offer / webrtc_answer / webrtc_ice:
payload := map[string]interface{}{"type": m.Type, "from": c.UserID, ...}
Hub.SendTo(m.To, payload)（未接続ならエラーを送信元へ webrtc_error 返信）
任意の未知 type は無視
３．通話要求offerを出すメソッド。
４．offerを受け取り、通知を出すメソッドは`useRecieveNOtification.ts`を利用できるか利用できるなら利用して実装。
５．answerを送信するメソッド。

## 次に実装したいもの
startchat/page.tsxにて¥会議を始めるボタンを押すと、`useStartChat.ts`のstartCallをトリガーを引き、answerが返ってくる。両者が通信完了したすなわち`pc.connectionState`がconnectedになったらそれがわかるようなUI loading/page.tsxに作成してほしいのだが
## 注意事項


##セキュリティ/運用
from はサーバ側で c.UserID を強制（なりすまし防止）
宛先が未接続ならエラーを送信元へ返却
メッセージサイズ・頻度（ICEスパム）にレート制限検討