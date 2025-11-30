import { fetchWithAuth } from "./fetchWithAuth";
import { API_BASE, WS_BASE } from "./apiBase";

export const startWebSocket = async(userID: number,token: string) => {

  await fetchWithAuth(`${API_BASE}/api/websocket/init`, {
    method:"POST",
    credentials: "include",
  });

  //websocket/initでtoken認証は済ませているからここではsession認証だけでいい
  const ws = new WebSocket(`${WS_BASE}/api/session/websocket`);

  ws.onopen = () => {
    console.log("✅ WebSocket接続成功");
    ws.send(JSON.stringify({ message: "Hello Server!" }));
  };

  ws.onerror = (error) => {
    console.error("❌ WebSocketエラー:", error);
  };

  ws.onclose = () => {
    console.log("🔌 WebSocket切断");
  };

  return ws;
};
