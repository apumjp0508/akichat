"use client";

import { useEffect, useState, useRef } from "react";

const mockUsers = [
  { id: 1, username: "太郎", avatar: "👨" },
  { id: 2, username: "花子", avatar: "👩" },
  { id: 3, username: "健一", avatar: "🧑" },
];

const mockMessages = [
  { id: 1, userId: 1, text: "こんにちは！" },
  { id: 2, userId: 2, text: "やっほー" },
  { id: 3, userId: 1, text: "会議の件どうする？" },
];

export default function ChatRoom() {
  const [messages, setMessages] = useState(mockMessages);
  const [newMessage, setNewMessage] = useState("");
  const [leftWidth, setLeftWidth] = useState(300);
  const resizerRef = useRef(null);
  const isResizing = useRef(false);

  useEffect(() => {
    const handleMouseMove = (e) => {
      if (!isResizing.current) return;
      const newWidth = e.clientX;
      if (newWidth > 200 && newWidth < window.innerWidth - 200) {
        setLeftWidth(newWidth);
      }
    };

    const handleMouseUp = () => {
      isResizing.current = false;
    };

    window.addEventListener("mousemove", handleMouseMove);
    window.addEventListener("mouseup", handleMouseUp);

    return () => {
      window.removeEventListener("mousemove", handleMouseMove);
      window.removeEventListener("mouseup", handleMouseUp);
    };
  }, []);

  const handleSend = () => {
    if (newMessage.trim() === "") return;
    const mockUser = mockUsers[0];
    setMessages([...messages, { id: Date.now(), userId: mockUser.id, text: newMessage }]);
    setNewMessage("");
  };

  return (
    <main className="flex h-screen w-full bg-gray-100">
      {/* 左：参加者欄 */}
      <div style={{ width: leftWidth }} className="bg-white border-r p-4 flex flex-col items-center">
        <h2 className="text-lg font-bold mb-4">参加者</h2>
        <div className="grid grid-cols-2 gap-4">
          {mockUsers.map((user) => (
            <div key={user.id} className="flex flex-col items-center">
              <div className="text-4xl">{user.avatar}</div>
              <div className="text-sm text-gray-700 mt-1">{user.username}</div>
            </div>
          ))}
        </div>
        <button className="mt-6 px-4 py-2 bg-red-500 text-white rounded-full hover:bg-red-600">
          退出
        </button>
      </div>

      {/* リサイズバー */}
      <div
        ref={resizerRef}
        onMouseDown={() => (isResizing.current = true)}
        className="w-2 cursor-col-resize bg-gray-300"
      ></div>

      {/* 右：DM欄 */}
      <div className="flex-1 flex flex-col p-4 bg-white">
        <div className="flex-1 overflow-y-auto space-y-4 p-2 rounded-lg shadow">
          {messages.map((msg) => {
            const user = mockUsers.find((u) => u.id === msg.userId);
            return (
              <div key={msg.id} className="flex items-start space-x-2">
                <div className="text-2xl">{user?.avatar}</div>
                <div>
                  <div className="text-sm font-semibold text-gray-700">{user?.username}</div>
                  <div className="text-sm bg-gray-200 rounded px-3 py-2 mt-1 inline-block">
                    {msg.text}
                  </div>
                </div>
              </div>
            );
          })}
        </div>

        {/* メッセージ入力欄 */}
        <div className="mt-4 flex items-center">
          <input
            type="text"
            className="flex-1 border border-gray-300 rounded-full px-4 py-2 focus:outline-none"
            placeholder="メッセージを入力..."
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSend()}
          />
          <button
            onClick={handleSend}
            className="ml-2 px-4 py-2 bg-blue-500 text-white rounded-full hover:bg-blue-600"
          >
            送信
          </button>
        </div>
      </div>
    </main>
  );
}