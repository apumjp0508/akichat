import { getWithAuth } from "./getwithAuth";
import { API_BASE } from "./apiBase";

export async function fetchFriends() {
  try {
    const data = await getWithAuth(`${API_BASE}/api/fetch/friends`);
    return data.friends;
  } catch (error) {
    console.error("Error fetching friends:", error);
    return [];
  }
}
