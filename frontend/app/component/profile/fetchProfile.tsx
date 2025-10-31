"use client";

import { useRouter } from "next/navigation";

export default function FetchProfilePage() {
  const router = useRouter();

    const fetchProfile = async () => {
        router.push("/profile");
  };
  return (
    <button
    onClick={fetchProfile}
    className="bg-purple-500 text-white py-2 px-4 rounded hover:bg-purple-600 transition"
    >
      プロフィールを確認する
    </button>
  );
}
