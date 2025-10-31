"use client";

import { useState, useEffect } from "react";
import { useUserStore } from "../../../../lib/store/userStore";
import { postWithAuth } from "../../../utils/PostWithAuth";
import { useReceiveNotification } from "../../../utils/RecieveNotification";

export default function NotificationListener({ userID }: { userID: number }) {
  const [notifications, setNotifications] = useState<Record<number, string>>({});
  const { user } = useUserStore();
  const token = user.token;

  useReceiveNotification(userID, token, (msg) => {
    if (msg.type === "friend_request") {
      const reqUserId = Number(msg.requestUserID);

      // 🔹 既存の通知オブジェクトに追加
      setNotifications((prev) => ({
        ...prev,
        [reqUserId]: msg.message,
      }));

      alert(`🔔 ${msg.message}`);
    }
  });

  const ApproveRequest = async (requestUserId,userID) =>{
    try {
      const res = await postWithAuth("http://localhost:8080/api/friend/request/approve",{
        requestUserId: requestUserId,
        userID: userID
      })

      if (res.ok) {
          const data = await res.json();
          console.log("フレンド申請成功:", data);
          alert("フレンド申請を送信しました。");
        } else {
          const errData = await res.json();
          console.error("フレンド申請失敗:", errData);
          alert(`フレンド申請に失敗: ${errData.error}`);
        }
      } catch (error) {
        console.error("Error during friend request:", error);
        alert("フレンド申請中にエラーが発生しました。");
      }
  }

  return (
    <div className="fixed bottom-4 right-4 space-y-2">
      {/* 🔹 Object.entries() で [key, value] に分けてループ */}
      {Object.entries(notifications).map(([reqId, text]) => (
        <div
          key={reqId}
          className="bg-blue-500 text-white px-4 py-2 rounded shadow-md animate-bounce"
        >
          <p>
            <strong>From User ID:</strong> {reqId}
          </p>
          <button
            type="submit"
            onClick={() => ApproveRequest(reqId, userID)}
          >
            リクエスト承認
          </button>
          <p>{text}</p>
        </div>
      ))}
    </div>
  );
}
