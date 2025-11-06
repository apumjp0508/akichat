"use client";

import { useEffect } from "react";
import FetchFriends  from "../../component/Friends/fetchFriends";
import { useState } from "react";
import { fetchFriends } from "../../utils/fetchFriends";
import { useRouter } from "next/navigation";

export default function ChatStartButton() {
    const [ friends, setFriends ] = useState<Array<{ ID: number; username: string }>>([]);
    const router = useRouter();
    useEffect(() => {
        const fetchData = async () => {
            const friends = await fetchFriends();
            setFriends(friends);
        };  
        fetchData();
    }, []);

    return (
        <main className="flex min-h-screen items-center justify-center bg-gradient-to-r from-blue-500 to-pink-400 text-white text-4xl font-bold">
            <div className="w-full max-w-md mx-auto bg-white rounded-2xl shadow-md p-6">
                <h1 className="text-2xl font-bold text-center text-pink-600 mb-6">🎉 友達リスト</h1>

                {friends.length === 0 ? (
                    <p className="text-center text-gray-500">友達はまだいません。</p>
                ) : (
                    <ul className="space-y-4">
                    {friends.map((friend) => (
                        <li
                        key={friend.ID}
                        className="flex items-center justify-between bg-pink-100 rounded-xl px-4 py-3 shadow-sm hover:bg-pink-200 transition"
                        >
                        <span className="text-pink-800 font-medium">{friend.username}</span>
                        <button
                            className="bg-pink-500 hover:bg-pink-600 text-white text-sm font-bold py-1 px-3 rounded-full shadow transition"
                            onClick={() => router.push("/loading")}
                        >
                            会議を始める    
                        </button>
                        </li>
                    ))}
                    </ul>
                )}
            </div>
        </main>
    );
}