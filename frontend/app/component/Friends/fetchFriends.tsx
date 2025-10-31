"use client";

import { useEffect, useState } from "react";
import { getWithAuth } from "../../utils/GetwithAuth";
import { useUserStore } from "../../../lib/store/userStore";

export default function FetchFriendPage() {
  const [friends, setFriends] = useState<Array<{ id: number; name: string }>>([]);
  const { user } = useUserStore();
  const token = user.token;

  useEffect(() => {
    fetchFriends();
  }, []);

    const fetchFriends = async () => {
    try {
      const res = await getWithAuth("http://localhost:8080/api/fetch/friends")

      if (res.ok) {
        const data = await res.json();
        setFriends(data.friends);
      } else {
        const data = await res.json();
        console.log(data);
        alert(`友達リスト取得失敗: ${data.error}`);
      }
    } catch (error) {
      console.error("Error fetching friends:", error);
      alert("友達リスト取得中にエラーが発生しました。");
    }
  };
  return (
    <div>
      <h2 className="text-xl font-bold mb-4">友達リスト</h2>
      <ul>
        {friends.map((friend) => (
          <li key={friend.id} className="mb-2">
            {friend.name}
          </li>
        ))}
      </ul>
    </div>
  );
}