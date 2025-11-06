import { getWithAuth } from "./GetwithAuth";

export async function fetchMe() {
  try {
    const data = await getWithAuth("http://localhost:8080/api/getMe");
    return data;
  } catch (error) {
    console.error("Error fetching me:", error);
    return null;
  }
}