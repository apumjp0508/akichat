
import { fetchWithAuth } from "./useFetchWithAuth";

export async function postWithAuth<T = any>(url: string, payload: any): Promise<T> {
  const res = await fetchWithAuth(url, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    const errorData = await res.json();
    throw new Error(errorData.error || "エラーが発生しました");
  }

  console.log("Response status:", res.status);
  return res.json();
}
