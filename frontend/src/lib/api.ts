const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001";

// Generic fetch wrapper with error handling
async function fetchApi<T>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  const url = `${API_URL}${endpoint}`;

  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.error || `HTTP error! status: ${response.status}`);
  }

  return response.json();
}

// User API
export const userApi = {
  create: (username: string, displayName: string) =>
    fetchApi<{ id: string; username: string; display_name: string }>(
      "/api/users",
      {
        method: "POST",
        body: JSON.stringify({ username, display_name: displayName }),
      }
    ),

  getById: (id: string) =>
    fetchApi<{ id: string; username: string; display_name: string }>(
      `/api/users/${id}`
    ),

  getByUsername: (username: string) =>
    fetchApi<{ id: string; username: string; display_name: string }>(
      `/api/users/username/${username}`
    ),
};

// Room API
export const roomApi = {
  list: () =>
    fetchApi<
      Array<{
        id: string;
        name: string;
        description: string | null;
        is_private: boolean;
        member_count: number;
        online_count: number;
      }>
    >("/api/rooms"),

  getById: (id: string) =>
    fetchApi<{
      id: string;
      name: string;
      description: string | null;
      is_private: boolean;
    }>(`/api/rooms/${id}`),

  create: (name: string, description?: string, isPrivate = false) =>
    fetchApi<{ id: string; name: string }>("/api/rooms", {
      method: "POST",
      body: JSON.stringify({ name, description, is_private: isPrivate }),
    }),

  join: (roomId: string, userId: string) =>
    fetchApi<{ message: string }>(`/api/rooms/${roomId}/join`, {
      method: "POST",
      body: JSON.stringify({ user_id: userId }),
    }),

  getMembers: (roomId: string) =>
    fetchApi<Array<{ id: string; username: string; display_name: string }>>(
      `/api/rooms/${roomId}/members`
    ),
};

// Message API
export const messageApi = {
  getByRoom: (roomId: string, limit = 50, offset = 0) =>
    fetchApi<{
      messages: Array<{
        id: string;
        room_id: string;
        user_id: string;
        content: string;
        message_type: string;
        created_at: string;
        username: string;
        display_name: string;
      }>;
      limit: number;
      offset: number;
    }>(`/api/rooms/${roomId}/messages?limit=${limit}&offset=${offset}`),
};

export default {
  user: userApi,
  room: roomApi,
  message: messageApi,
};
