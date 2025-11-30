"use client";

import { useState } from "react";
import { useUserStore } from "../../../lib/store/userStore";
import { postWithAuth } from "../../utils/postWithAuth";
import { API_BASE } from "../../utils/apiBase";
import { useReceiveNotification } from "../../hook/useRecieveNotification";
import { acceptCall } from "../../utils/startChatOffer";
import { handleRemoteAnswer, applyRemoteIce } from "../../utils/callSession";

export default function NotificationListener({ userID }: { userID: number }) {
  const [notifications, setNotifications] = useState<Record<number, string>>({});
  const { user } = useUserStore();
  const token = user.token;
  const peerByUser = new Map<number, RTCPeerConnection>();

  useReceiveNotification(
    userID,
    token,
    (msg) => {
      if (msg.type === "friend_request") {
        const reqUserID = Number(msg.requestUserID);

        // 🔹 既存の通知オブジェクトに追加
        setNotifications((prev) => ({
          ...prev,
          [reqUserID]: msg.message,
        }));

        alert(`🔔 ${msg.message}`);
      }
    },
    {
      // オファー受信: 承諾でアンサーを返す
      onOffer: async (ws, { from, sdp }) => {
        try {
          const accept = window.confirm(`📞 User ${from} からの通話リクエスト。受けますか？`);
          if (!accept) return;
          // 受信用 PeerConnection を作成し、アンサー送信までをユーティリティに委譲
          const pc = await acceptCall(ws, from, sdp);
          peerByUser.set(from, pc);
        } catch (e) {
          console.error("onOffer handling failed:", e);
          alert("通話接続中にエラーが発生しました。");
        }
      },
      // 呼び出し側: アンサーを適用
      onAnswer: async (_ws, { from, sdp }) => {
        try {
          await handleRemoteAnswer(from, sdp);
        } catch (e) {
          console.error("apply remote answer failed:", e);
        }
      },
      // 相手のICE候補を適用（統一関数へ委譲）
      onIce: async (_ws, { from, candidate }) => {
        try {
          // Mapがある場合はMapを、なければ発信側（単一セッション）に適用
          const pc = peerByUser.get(from);
          await applyRemoteIce(from, candidate, pc ? peerByUser : undefined);
        } catch (e) {
          console.error("addIceCandidate failed:", e);
        }
      },
    }
  );

  const ApproveRequest = async (requestUserID,userID) =>{
    try {
      const data = await postWithAuth(`${API_BASE}/api/friend/request/approve`,{
        requestUserID: Number(requestUserID),
        userID: Number(userID),
      })

      console.log("フレンド承認成功:", data);
      alert("フレンド申請を承認しました。");
      
      } catch (error) {
        console.log("Approving friend request from user ID:", requestUserID);
        console.log("Current user ID:", userID);
        console.error("Error during friend request:", error);
        alert("フレンド申請中にエラーが発生しました。");
      }
  }

  return (
    <div className="fixed bottom-4 right-4 space-y-2">
      {/* 🔹 Object.entries() で [key, value] に分けてループ */}
      {Object.entries(notifications).map(([reqID, text]) => (
        <div
          key={reqID}
          className="bg-blue-500 text-white px-4 py-2 rounded shadow-md animate-bounce"
        >
          <p>
            <strong>From User ID:</strong> {reqID}
          </p>
          <button
            type="submit"
            onClick={() => ApproveRequest(reqID, userID)}
          >
            リクエスト承認
          </button>
          <p>{text}</p>
        </div>
      ))}
    </div>
  );
}
