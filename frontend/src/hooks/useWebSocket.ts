"use client";

import { useState, useEffect, useRef, useCallback } from "react";
import { Message, OnlineUser, TypingUser, WSMessage, WSMessageType } from "@/types";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://127.0.0.1:3001";
const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://127.0.0.1:3001";

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

  const historyReceivedRef = useRef(false);
  const historyFetchAttemptedRef = useRef(false);
  const historyProbeTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const messagesLenRef = useRef(0);
  useEffect(() => {
    messagesLenRef.current = messages.length;
  }, [messages.length]);

  // Reset per-room state when room changes
  useEffect(() => {
    setMessages([]);
    setOnlineUsers([]);
    setTypingUsers([]);
    setIsConnected(false);
    setError(null);
    historyReceivedRef.current = false;
    historyFetchAttemptedRef.current = false;
    if (historyProbeTimeoutRef.current) {
      clearTimeout(historyProbeTimeoutRef.current);
      historyProbeTimeoutRef.current = null;
    }
  }, [roomId]);

  const fetchHistoryHttp = useCallback(async () => {
    if (!roomId) return;
    if (historyFetchAttemptedRef.current) return;
    historyFetchAttemptedRef.current = true;

    try {
      const res = await fetch(`${API_URL}/api/rooms/${roomId}/messages?limit=50&offset=0`);
      if (!res.ok) {
        return;
      }
      const data = await res.json();
      const next = (data?.messages as Message[]) || [];
      // Only set if we still haven't received WS history
      if (!historyReceivedRef.current) {
        setMessages(next);
      }
    } catch {
      // Ignore: this is a best-effort fallback
    }
  }, [roomId]);

  // Connect to WebSocket
  const connect = useCallback(() => {
    // Guard against empty parameters
    if (!roomId || !userId || userId === "") {
      console.log("⏸️ WebSocket: Waiting for roomId and userId...");
      return;
    }

    // Close existing connection first
    if (wsRef.current) {
      console.log("🔄 Closing existing WebSocket connection...");
      wsRef.current.onclose = null; // Prevent triggering reconnect
      wsRef.current.close();
      wsRef.current = null;
    }

    // Clear any pending reconnect
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }

    const wsUrl = `${WS_URL}/ws/${roomId}?userId=${userId}&username=${encodeURIComponent(username)}&displayName=${encodeURIComponent(displayName)}`;

    console.log("🔗 Connecting to WebSocket:", wsUrl);

    try {
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;

      ws.onopen = () => {
        console.log("✅ WebSocket connected");
        setIsConnected(true);
        setError(null);
        reconnectAttempts.current = 0;

        // If WS doesn't deliver history quickly (or ever), fall back to HTTP.
        if (historyProbeTimeoutRef.current) {
          clearTimeout(historyProbeTimeoutRef.current);
        }
        historyProbeTimeoutRef.current = setTimeout(() => {
          if (!historyReceivedRef.current && messagesLenRef.current === 0) {
            fetchHistoryHttp();
          }
        }, 1200);
      };

      ws.onclose = (event) => {
        console.log("❌ WebSocket disconnected:", event.code, event.reason);
        setIsConnected(false);

        // Still allow loading history even if WS failed.
        if (!historyReceivedRef.current && messagesLenRef.current === 0) {
          fetchHistoryHttp();
        }

        // Only reconnect if this is still the current connection
        if (wsRef.current === ws) {
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
        }
      };

      ws.onerror = (event) => {
        // Next.js dev overlay treats console.error as an app error; use warn instead.
        const readyState = ws.readyState;
        console.warn("⚠️ WebSocket error", {
          readyState,
          url: wsUrl,
          type: (event as Event)?.type,
        });
        // Only set error if this is still the current connection
        if (wsRef.current === ws) {
          setError("เกิดข้อผิดพลาดในการเชื่อมต่อ WebSocket");
        }
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
      fetchHistoryHttp();
    }
  }, [roomId, userId, username, displayName]);

  // Handle incoming messages
  const handleMessage = useCallback((data: WSMessage) => {
    switch (data.type) {
      case "message":
        if (data.payload) {
          setMessages((prev) => [...prev, data.payload as Message]);
        }
        break;

      case "history":
        historyReceivedRef.current = true;
        setMessages((data.payload as Message[]) || []);
        break;

      case "online_users":
        setOnlineUsers((data.payload as OnlineUser[]) || []);
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
      if (historyProbeTimeoutRef.current) {
        clearTimeout(historyProbeTimeoutRef.current);
        historyProbeTimeoutRef.current = null;
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
