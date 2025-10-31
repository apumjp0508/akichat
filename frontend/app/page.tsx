"use client";

import { useState,useEffect } from "react";
import { useUserStore } from "../lib/store/userStore";
import { fetchConnectedUsers } from "./utils/fetchConnectedUsers";
import Notification from "./component/home/Notification/notification";
import AuthDashboard from "./component/authentification/dashboard";
import GuestDashboard from "./component/authentification/guestDashBoard";

export default function HomePage() {
  const [isLogin, setIsLogin] = useState(true);
  const [connectedUsers, setConnectedUsers] = useState<number[]>([]);
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
      (async () => {
        const data = await fetchConnectedUsers();
        if (data?.connected_users) {
          setConnectedUsers(data.connected_users);
        }
      })();
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
      <h1>こんにちは{userID}さん！</h1>
      <Notification userID={userID}/>
      {
        isLogin ? <AuthDashboard/> : <GuestDashboard/>
      }
      <h2>接続中のユーザー一覧{connectedUsers}</h2>
    </div>
  );
}
