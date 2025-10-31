"use client";

import { useEffect } from "react";
import { startWebSocket } from "./StartWebSocket";

export function useReceiveNotification(
  userID: number,
  token: string,
  onMessage: (msg: any) => void
) {
  useEffect(() => {
  if (!userID || !token) return;

  let socket: WebSocket | null = null;

  const initWebSocket = async () => {
    socket = await startWebSocket(userID, token); // ✅ awaitでPromiseを解決

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log("📩 Received:", data);
        onMessage(data);
      } catch (err) {
        console.error("JSON parse error:", err);
      }
    };
  };

  initWebSocket();

  return () => {
    socket?.close(); // ✅ 安全に呼び出し
  };
}, [userID, token, onMessage]);

}
