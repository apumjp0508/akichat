"use client";

import { useRouter } from "next/navigation";

export default function GetMePage() {
  const router = useRouter();

    const getMe = async () => {
        router.push("/getMe");
  };
  return (
    <button
    onClick={getMe}
    className="bg-purple-500 text-white py-2 px-4 rounded hover:bg-purple-600 transition"
    >
      プロフィールを確認する
    </button>
  );
}
