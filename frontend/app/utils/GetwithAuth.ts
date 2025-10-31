import { fetchWithAuth } from "./fetchWithAuth";

export async function getWithAuth<T = any>(url: string): Promise<T> {
  const res = await fetchWithAuth(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!res.ok) {
    const errorData = await res.json();
    throw new Error(errorData.error || "エラーが発生しました");
  }

  return res.json();
}