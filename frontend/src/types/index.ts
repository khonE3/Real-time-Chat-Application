// TypeScript interfaces for the chat application

// User types
export interface User {
  id: string;
  username: string;
  display_name: string;
  avatar_url?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateUserRequest {
  username: string;
  display_name: string;
  avatar_url?: string;
}

// Room types
export interface Room {
  id: string;
  name: string;
  description: string | null;
  is_private: boolean;
  created_by: string | null;
  created_at: string;
}

export interface RoomWithMembers extends Room {
  member_count: number;
  online_count: number;
  unread_count?: number;
}

export interface CreateRoomRequest {
  name: string;
  description?: string;
  is_private?: boolean;
}

// Message types
export type MessageType = "text" | "system" | "typing" | "presence";

export interface Message {
  id: string;
  room_id: string;
  user_id: string | null;
  content: string;
  message_type: MessageType;
  created_at: string;
  username?: string;
  display_name?: string;
  avatar_url?: string;
}

export interface SendMessageRequest {
  content: string;
}

// WebSocket message types
export type WSMessageType =
  | "message"
  | "typing"
  | "stop_typing"
  | "presence"
  | "history"
  | "online_users"
  | "error"
  | "join"
  | "leave";

export interface WSMessage<T = unknown> {
  type: WSMessageType;
  payload: T;
}

export interface WSIncomingMessage {
  type: WSMessageType;
  content?: string;
  user_id?: string;
}

// Typing indicator
export interface TypingUser {
  user_id: string;
  username: string;
  display_name: string;
  is_typing: boolean;
}

// Online user
export interface OnlineUser {
  user_id: string;
  username: string;
  display_name: string;
  avatar_url?: string;
  last_seen: string;
}

// Presence update
export interface PresencePayload {
  user_id: string;
  username: string;
  display_name: string;
  is_online: boolean;
}

// API Response types
export interface ApiResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  limit: number;
  offset: number;
  total?: number;
}

// Chat state
export interface ChatState {
  messages: Message[];
  onlineUsers: OnlineUser[];
  typingUsers: TypingUser[];
  isConnected: boolean;
  error: string | null;
}
