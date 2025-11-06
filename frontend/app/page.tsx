"use client";

import { useState,useEffect } from "react";
import { useUserStore } from "../lib/store/userStore";
import Notification from "./component/home/approveNotification";
import AuthDashboard from "./component/AuthDashboard/page";
import ConnectedUsers from "./component/home/connectedUser";
import GuestDashboard from "./component/GuestDashboard/page";

export default function HomePage() {
  const [isLogin, setIsLogin] = useState(true);
  const { user } = useUserStore();
  const token = user.token;
  const userID = user.id;

  const isTokenExpired = (token: string): boolean => {
    try {
      const payload = JSON.parse(atob(token.split(".")[1]));
      const currentTime = Math.floor(Date.now() / 1000);
      return payload.exp < currentTime;
    } catch (e) {
      console.error("Invalid token:", e);
      return true;
    }
  };

  useEffect(() => {
    
    if (!token) {
      setIsLogin(false);
      return;
    }

    if (isTokenExpired(token)) {
      setIsLogin(false);
    } else {
      setIsLogin(true);
    }
  }, [token]);
  return (
    <div>
      <main className="flex min-h-screen items-center justify-center bg-gradient-to-r from-blue-500 to-pink-400 text-white text-4xl font-bold">
        <Notification userID={userID}/>
        {
          isLogin ? (
            <>
              <h1 className="text-4xl font-bold">こんにちは{userID}さん！</h1>
              <AuthDashboard/>
              <ConnectedUsers/>
            </> 
          ): (
            <GuestDashboard/>
          )
        }
      </main>
    </div>
  );
}
