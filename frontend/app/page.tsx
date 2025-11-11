"use client";

import { useState, useEffect, use } from "react";
import { useUserStore } from "../lib/store/userStore";
import Notification from "./component/home/approveNotification";
import AuthDashboard from "./component/AuthDashboard/page";
import GuestDashboard from "./component/GuestDashboard/page";
import { checkCookie } from "./utils/useCheckCookie";

export default function HomePage() {
  const [isChecking, setIsChecking] = useState(true);
  const [isLogin, setIsLogin] = useState(false);
  const [ userID, setUserID ] = useState<number | null>(null);
  const [ userName, setUserName ] = useState<string | null>(null);
  const [ userEmail, setUserEmail ] = useState<string | null>(null);

  useEffect(() => {
    // ✅ 初回マウント時にログイン状態を確認
    const verifyLogin = async () => {
      setIsChecking(true);
      const result = await checkCookie();

      if (result.loggedIn) {
        setIsLogin(true);
        setUserID(result.id);
        setUserName(result.name);
        setUserEmail(result.email);
        useUserStore.getState().setUser({
            id: result.id,
            name: result.name,
            email: result.email,
        });
        console.log("✅ ログイン済み:", result.message);
      } else {
        setIsLogin(false);
        console.log("🚪 未ログイン:", result.message);
      }
      setIsChecking(false);
    };

    verifyLogin();
  }, []);

  if (isChecking) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-100 text-gray-800 text-2xl">
        認証確認中...
      </div>
    );
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-gradient-to-r from-blue-500 to-pink-400 text-white text-4xl font-bold">
      {isLogin ? (
        <>
          <Notification userID={userID} />
          <div className="text-center">
            <h1 className="text-4xl font-bold mb-6">
              こんにちは {userID ?? "ゲスト"} さん！
            </h1>
            <h2>
              ユーザー名: {userName}<br />
              メール: {userEmail ?? ""}
            </h2>
            <AuthDashboard />
          </div>
        </>
      ) : (
        <GuestDashboard />
      )}
    </main>
  );
}
