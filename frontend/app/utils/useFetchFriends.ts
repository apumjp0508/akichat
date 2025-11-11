import { getWithAuth } from "./useGetwithAuth";

export async function fetchFriends() {
  try {
    const data = await getWithAuth("http://localhost:8080/api/fetch/friends");
    return data.friends;
  } catch (error) {
    console.error("Error fetching friends:", error);
    return [];
  }
}
