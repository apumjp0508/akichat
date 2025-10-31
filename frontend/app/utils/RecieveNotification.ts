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

    const socket = startWebSocket(userID, token);

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log("📩 Received:", data);
        onMessage(data);
      } catch (err) {
        console.error("JSON parse error:", err);
      }
    };

    return () => {
      socket.close();
    };
  }, [userID, token, onMessage]);
}
