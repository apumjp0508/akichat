import { fetchWithAuth } from "./fetchWithAuth";

export const startWebSocket = async(userID: number,token: string) => {

  await fetchWithAuth("http://localhost:8080/api/websocket/init", {
    method:"POST",
    credentials: "include",
  });

  //websocket/initでtoken認証は済ませているからここではsession認証だけでいい
  const ws = new WebSocket("ws://localhost:8080/api/session/websocket");

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
