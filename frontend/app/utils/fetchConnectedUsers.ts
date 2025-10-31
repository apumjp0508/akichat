import { getWithAuth } from "./GetwithAuth";

export async function fetchConnectedUsers() {
  try {
    const data = await getWithAuth("http://localhost:8080/api/connected-users");

    return data;
  } catch (error) {
    console.error("Error fetching connectedUsers:", error);
    return null;
  }
}
