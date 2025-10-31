import { fetchWithAuth } from "./fetchWithAuth";

export const startWebSocket = async(userID: number,token: string) => {

  await fetchWithAuth("/api/websocket/init", {
    method:"POST",
    credentials: "include",
  });

  const ws = new WebSocket("ws://localhost:8080/api/websocket");

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
