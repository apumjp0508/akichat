"use client";
import { useRouter } from "next/navigation";

export default function Loading() {
    const router = useRouter();
    return (
      <main className="flex min-h-screen items-center justify-center bg-gradient-to-r from-blue-500 to-pink-400 text-white">
        <div className="w-full max-w-md mx-auto bg-white rounded-2xl shadow-lg p-8 flex flex-col items-center">
          <div className="w-12 h-12 border-4 border-blue-400 border-t-transparent rounded-full animate-spin mb-4"></div>
          <h1 className="text-xl font-semibold text-blue-600 mb-2">💬 チャットルームへ接続中</h1>
          <p className="text-gray-600">サーバーとの接続を確立しています...</p>
          <button
          onClick={() => router.push("/chatRoom")}
          className="bg-blue-500 text-white py-2 px-4 my-2 rounded-lg hover:bg-blue-600 transition"
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