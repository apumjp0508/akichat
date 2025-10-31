import { useUserStore } from "../../lib/store/userStore";
let refreshPromise: Promise<boolean> | null = null;


export async function fetchWithAuth(input: RequestInfo, init: RequestInit = {}) {
  const token = useUserStore.getState().user.token || localStorage.getItem('token');
  const headers = new Headers(init.headers || {});
  if (token) headers.set('Authorization', `Bearer ${token}`);
  const res = await fetch(input, { ...init, headers, credentials: 'include' });
  if (res.status !== 401) return res;

  if (!refreshPromise) {
    refreshPromise = (async () => {
      try {
        const r = await fetch('/api/auth/refresh', { method: 'POST', credentials: 'include' });
        refreshPromise = null;
        return r.ok;
      } catch {
        refreshPromise = null;
        return false;
      }
    })();
  }
  const refreshed = await refreshPromise;
  if (!refreshed) {
    // logout
    useUserStore.getState().clearUser();
    throw new Error('Session expired');
  }
  // refresh 成功 => 再試行
  const token2 = useUserStore.getState().user.token || localStorage.getItem('token');
  if (token2) headers.set('Authorization', `Bearer ${token2}`);
  return fetch(input, { ...init, headers, credentials: 'include' });
}