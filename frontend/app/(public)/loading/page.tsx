"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { subscribeConnectionState, getConnectionState } from "../../utils/callSession";

export default function Loading() {
    const router = useRouter();
    const [state, setState] = useState<string>(() => getConnectionState());

    useEffect(() => {
        const unsub = subscribeConnectionState((s) => setState(s));
        return () => unsub();
    }, []);

    const isConnected = state === "connected";

    return (
      <main className="flex min-h-screen items-center justify-center bg-gradient-to-r from-blue-500 to-pink-400 text-white">
        <div className="w-full max-w-md mx-auto bg-white rounded-2xl shadow-lg p-8 flex flex-col items-center">
          {!isConnected ? (
            <>
              <div className="w-12 h-12 border-4 border-blue-400 border-t-transparent rounded-full animate-spin mb-4"></div>
              <h1 className="text-xl font-semibold text-blue-600 mb-2">💬 ピア接続中...</h1>
              <p className="text-gray-600">WebRTC 状態: {state}</p>
            </>
          ) : (
            <>
              <div className="w-12 h-12 bg-green-400 rounded-full mb-4 shadow-[0_0_12px_#22c55e]"></div>
              <h1 className="text-xl font-semibold text-green-600 mb-2">✅ 接続完了</h1>
              <p className="text-gray-600">WebRTC 状態: {state}</p>
            </>
          )}
          <button
            onClick={() => router.push("/chatRoom")}
            className={`py-2 px-4 my-4 rounded-lg transition ${isConnected ? "bg-blue-500 hover:bg-blue-600 text-white" : "bg-gray-300 text-gray-500 cursor-not-allowed"}`}
            disabled={!isConnected}
          >
            会議を始める
          </button>
          <button
            onClick={() => router.back()}
            className="bg-blue-500 text-white py-2 px-4 rounded-lg hover:bg-blue-600 transition"
          >
            戻る
          </button>
        </div>
      </main>
    );
}