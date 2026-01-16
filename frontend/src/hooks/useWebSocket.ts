"use client";

import { useState, useEffect, useRef, useCallback } from "react";
import { Message, OnlineUser, TypingUser, WSMessage, WSMessageType } from "@/types";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:3001";

interface UseWebSocketReturn {
  messages: Message[];
  onlineUsers: OnlineUser[];
  typingUsers: TypingUser[];
  isConnected: boolean;
  error: string | null;
  sendMessage: (content: string) => void;
  sendTyping: () => void;
  sendStopTyping: () => void;
}

export function useWebSocket(
  roomId: string,
  userId: string,
  username: string,
  displayName: string
): UseWebSocketReturn {
  const [messages, setMessages] = useState<Message[]>([]);
  const [onlineUsers, setOnlineUsers] = useState<OnlineUser[]>([]);
  const [typingUsers, setTypingUsers] = useState<TypingUser[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttempts = useRef(0);
  const maxReconnectAttempts = 5;

  // Connect to WebSocket
  const connect = useCallback(() => {
    if (!roomId || !userId) return;

    const wsUrl = `${WS_URL}/ws/${roomId}?userId=${userId}&username=${username}&displayName=${encodeURIComponent(displayName)}`;

    try {
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;

      ws.onopen = () => {
        console.log("✅ WebSocket connected");
        setIsConnected(true);
        setError(null);
        reconnectAttempts.current = 0;
      };

      ws.onclose = (event) => {
        console.log("❌ WebSocket disconnected:", event.code, event.reason);
        setIsConnected(false);

        // Attempt to reconnect
        if (reconnectAttempts.current < maxReconnectAttempts) {
          reconnectAttempts.current++;
          const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.current), 30000);
          console.log(`🔄 Reconnecting in ${delay}ms (attempt ${reconnectAttempts.current})`);

          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, delay);
        } else {
          setError("ไม่สามารถเชื่อมต่อได้ กรุณารีเฟรชหน้า");
        }
      };

      ws.onerror = (event) => {
        console.error("⚠️ WebSocket error:", event);
        setError("เกิดข้อผิดพลาดในการเชื่อมต่อ");
      };

      ws.onmessage = (event) => {
        try {
          const data: WSMessage = JSON.parse(event.data);
          handleMessage(data);
        } catch (e) {
          console.error("Failed to parse message:", e);
        }
      };
    } catch (e) {
      console.error("Failed to create WebSocket:", e);
      setError("ไม่สามารถสร้างการเชื่อมต่อได้");
    }
  }, [roomId, userId, username, displayName]);

  // Handle incoming messages
  const handleMessage = useCallback((data: WSMessage) => {
    switch (data.type) {
      case "message":
        setMessages((prev) => [...prev, data.payload as Message]);
        break;

      case "history":
        setMessages(data.payload as Message[]);
        break;

      case "online_users":
        setOnlineUsers(data.payload as OnlineUser[]);
        break;

      case "typing":
        const typingUser = data.payload as TypingUser;
        setTypingUsers((prev) => {
          const exists = prev.some((u) => u.user_id === typingUser.user_id);
          if (exists) return prev;
          return [...prev, typingUser];
        });
        break;

      case "stop_typing":
        const stopUser = data.payload as TypingUser;
        setTypingUsers((prev) =>
          prev.filter((u) => u.user_id !== stopUser.user_id)
        );
        break;

      case "presence":
        const presence = data.payload as {
          user_id: string;
          username: string;
          display_name: string;
          is_online: boolean;
        };
        if (presence.is_online) {
          setOnlineUsers((prev) => {
            const exists = prev.some((u) => u.user_id === presence.user_id);
            if (exists) return prev;
            return [
              ...prev,
              {
                user_id: presence.user_id,
                username: presence.username,
                display_name: presence.display_name,
                last_seen: new Date().toISOString(),
              },
            ];
          });
        } else {
          setOnlineUsers((prev) =>
            prev.filter((u) => u.user_id !== presence.user_id)
          );
          setTypingUsers((prev) =>
            prev.filter((u) => u.user_id !== presence.user_id)
          );
        }
        break;

      case "error":
        setError(data.payload as string);
        break;

      default:
        console.log("Unknown message type:", data.type);
    }
  }, []);

  // Send message
  const sendMessage = useCallback((content: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({
          type: "message" as WSMessageType,
          content,
        })
      );
    }
  }, []);

  // Send typing indicator
  const sendTyping = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({
          type: "typing" as WSMessageType,
        })
      );
    }
  }, []);

  // Send stop typing indicator
  const sendStopTyping = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({
          type: "stop_typing" as WSMessageType,
        })
      );
    }
  }, []);

  // Connect on mount
  useEffect(() => {
    connect();

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [connect]);

  return {
    messages,
    onlineUsers,
    typingUsers,
    isConnected,
    error,
    sendMessage,
    sendTyping,
    sendStopTyping,
  };
}
