"use client";

import { useRouter } from "next/navigation";
import GetMePage from "../profile/getMe"; 

export default function AuthDashboard() { 

    const router = useRouter();

    return (        
     <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-blue-100 via-pink-100 to-yellow-100 p-10">

      <div className="flex flex-col gap-5 w-72">
        <button
          onClick={() => router.push("/login")}
          className="bg-blue-500 text-white py-3 rounded-xl font-bold shadow-md hover:scale-105 hover:bg-blue-600 transition-transform duration-300"
        >
          🔐 ログイン画面へ
        </button>

        <button
          onClick={() => router.push("/register")}
          className="bg-green-500 text-white py-3 rounded-xl font-bold shadow-md hover:scale-105 hover:bg-green-600 transition-transform duration-300"
        >
          ✨ 新規登録画面へ
        </button>

        <div className="bg-white/60 rounded-xl p-4 shadow-inner backdrop-blur-sm">
          <GetMePage />
        </div>

        <button
          onClick={() => router.push("/searchFriends")}
          className="bg-pink-500 text-white py-3 rounded-xl font-bold shadow-md hover:scale-105 hover:bg-pink-600 transition-transform duration-300"
        >
          💬 友達一覧を見る
        </button>

        <button
          onClick={() => router.push("/searchNotFriend")}
          className="bg-yellow-400 text-gray-800 py-3 rounded-xl font-bold shadow-md hover:scale-105 hover:bg-yellow-500 transition-transform duration-300"
        >
          🔍 友達ではないユーザーを探す
        </button>

        <button
          onClick={() => router.push("/startChat")}
          className="bg-red-500 text-white py-3 rounded-xl font-bold shadow-md hover:scale-105 hover:bg-red-600 transition-transform duration-300"
        >
          💬 チャットを始める
        </button>
      </div>
    </div>
    );
}   