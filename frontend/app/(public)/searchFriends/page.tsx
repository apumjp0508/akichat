"use client";

import FetchFriendPage from "../../component/Friends/fetchFriends";

export default function SearchFriendsPage() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-10">
      <h1 className="text-2xl font-bold mb-6">友達を探す</h1>
      <div className="w-full max-w-md bg-white p-6 rounded shadow">
        <FetchFriendPage />
      </div>
    </div>
  );
}