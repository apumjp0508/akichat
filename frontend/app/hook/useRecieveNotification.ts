"use client";

import { useEffect } from "react";
import { startWebSocket } from "../utils/startWebSocket";
import { attachSignalingHandlers, SignalingHandlers } from "../utils/signaling";

export function useReceiveNotification(
  userID: number,
  token: string,
  onMessage: (msg: any) => void,
  signalingHandlers?: {
    //「onOffer は任意の関数であり、引数として (ws, data) を受け取り、何も返さない（返り値を使わない） 関数」
    onOffer?: (ws: WebSocket, data: { from: number; sdp: RTCSessionDescriptionInit }) => void;
    onAnswer?: (ws: WebSocket, data: { from: number; sdp: RTCSessionDescriptionInit }) => void;
    onIce?: (ws: WebSocket, data: { from: number; candidate: RTCIceCandidateInit }) => void;
    onError?: (ws: WebSocket, data: { reason: string; to?: number; originalType?: string }) => void;
  }
) {
  useEffect(() => {
    if (!userID || !token) return;

    let socket: WebSocket | null = null;
    let detachSignaling: (() => void) | null = null;

    const initWebSocket = async () => {
      socket = await startWebSocket(userID, token); // ✅ awaitでPromiseを解決

      // シグナリング受信（追加ハンドラ）
      const handlers: SignalingHandlers = {
        onOffer: (d) => signalingHandlers?.onOffer?.(socket!, d),
        onAnswer: (d) => signalingHandlers?.onAnswer?.(socket!, d),
        onIce: (d) => signalingHandlers?.onIce?.(socket!, d),
        onError: (d) => signalingHandlers?.onError?.(socket!, d),
      };
      detachSignaling = attachSignalingHandlers(socket, handlers, { onMessage });
    };

    initWebSocket();

    return () => {
      try {
        detachSignaling?.();
      } finally {
        socket?.close(); // ✅ 安全に呼び出し
      }
    };
  }, [userID, token, onMessage, signalingHandlers]);
}
