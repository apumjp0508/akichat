"use client";

import { useRouter } from "next/navigation";

export default function AuthDashboard() { 

    const router = useRouter();

    const handleGoToLogin = () => {
        router.push("/login");
    };

    const handleGoToRegister = () => {
        router.push("/register");
    };


    return (        
    <div className="p-10 flex flex-col items-center justify-center min-h-screen">
      <h1 className="text-2xl font-bold mb-6">ゲストダッシュボード</h1>

      <div className="flex flex-col gap-4 w-60">
        <button
          onClick={handleGoToLogin}
          className="bg-blue-500 text-white py-2 rounded hover:bg-blue-600 transition"
        >
          ログイン画面へ
        </button>

        <button
          onClick={handleGoToRegister}
          className="bg-green-500 text-white py-2 rounded hover:bg-green-600 transition"
        >
          新規登録画面へ
        </button>

      </div>
    </div>
    );
}   