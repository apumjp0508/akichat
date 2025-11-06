import { useUserStore } from "../../lib/store/userStore";
import { refreshTokenIfNeeded } from "./refreshTokenIfNeeded";

export async function fetchWithAuth(input: RequestInfo, init: RequestInit = {}) {
  const token = useUserStore.getState().user.token || localStorage.getItem("token");
  const headers = new Headers(init.headers || {});
  if (token) headers.set("Authorization", `Bearer ${token}`);

  let res = await fetch(input, { ...init, headers, credentials: "include" });

  if (res.status !== 401) return res;

  const refreshed = await refreshTokenIfNeeded();
  if (!refreshed) {
    useUserStore.getState().clearUser();
    throw new Error("Session expired");
  }

  const newToken = useUserStore.getState().user.token || localStorage.getItem("token");
  if (newToken) headers.set("Authorization", `Bearer ${newToken}`);

  // 再試行
  return fetch(input, { ...init, headers, credentials: "include" });
}
