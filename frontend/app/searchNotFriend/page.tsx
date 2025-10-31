"use client";

import { useState } from "react";
import { fetchWithAuth } from "../utils/fetchWithAuth";
import { postWithAuth } from "../utils/PostWithAuth";

export default function SearchNotFriendPage() {
  const [keyword, setKeyword] = useState("");
  const [userList, setUserList] = useState<Array<{ ID: number; username: string; email: string }>>([]);

  const handleSearchNotFriends = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const res = await postWithAuth("http://localhost:8080/api/search/notfriends",{
        keyword,
      });

      if (res.ok) {
        const data = await res.json();
        console.log("検索成功:", data);
        setUserList([data.user]); // 1件ならこう
      } else {
        const errData = await res.json();
        console.error("検索失敗:", errData);
        alert(`友達ではないユーザーの取得に失敗: ${errData.error}`);
      }
    } catch (error) {
      console.error("Error during fetching not friends:", error);
      alert("友達ではないユーザーの取得中にエラーが発生しました。");
    }
  };

  const FriendRequest = async (userID: number) => {
    try {
      const res = await postWithAuth("http://localhost:8080/api/friend/request",{
        to_user_id: userID ,
      });
      if (res.ok) {
          const data = await res.json();
          console.log("フレンド申請成功:", data);
          alert("フレンド申請を送信しました。");
        } else {
          const errData = await res.json();
          console.error("フレンド申請失敗:", errData);
          alert(`フレンド申請に失敗: ${errData.error}`);
        }
      } catch (error) {
        console.error("Error during friend request:", error);
        alert("フレンド申請中にエラーが発生しました。");
      }
    };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-10">
      <h1 className="text-2xl font-bold mb-6">友達ではないユーザーを探す</h1>
      <div className="w-full max-w-md bg-white p-6 rounded shadow">
        <form onSubmit={handleSearchNotFriends}>
          <input
            type="text"
            placeholder="ユーザー名で検索"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            className="border p-2 rounded w-full mb-4"
            required
          />
          <button
            type="submit"
            className="bg-blue-500 text-white p-2 rounded w-full hover:bg-blue-600"
          >
            検索
          </button>
        </form>

        {/* 検索結果を表示 */}
        <ul className="mt-4">
          {userList.map((u) => (
            <li key={u.ID} className="border-b py-2">
              <p>👤 {u.username}</p>
              <p className="text-sm text-gray-600">{u.email}</p>
              <button
                type="submit"
                onClick={() => FriendRequest(u.ID)}
                className="mt-2 bg-green-500 text-white p-1 rounded hover:bg-green-600"
              >
                フレンド申請
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
