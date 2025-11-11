"use client";

import { useState } from "react";
import { useUserStore } from "../../../lib/store/userStore";
import { postWithAuth } from "../../utils/usePostWithAuth";
import { useReceiveNotification } from "../../utils/useRecieveNotification";

export default function NotificationListener({ userID }: { userID: number }) {
  const [notifications, setNotifications] = useState<Record<number, string>>({});
  const { user } = useUserStore();
  const token = user.token;

  useReceiveNotification(userID, token, (msg) => {
    if (msg.type === "friend_request") {
      const reqUserID = Number(msg.requestUserID);

      // 🔹 既存の通知オブジェクトに追加
      setNotifications((prev) => ({
        ...prev,
        [reqUserID]: msg.message,
      }));

      alert(`🔔 ${msg.message}`);
    }
  });

  const ApproveRequest = async (requestUserID,userID) =>{
    try {
      const data = await postWithAuth("http://localhost:8080/api/friend/request/approve",{
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
