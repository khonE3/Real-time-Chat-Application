"use client";

import { useState, useEffect, useRef, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import { useWebSocket } from "@/hooks/useWebSocket";
import MessageList from "@/components/chat/MessageList";
import MessageInput from "@/components/chat/MessageInput";
import OnlineUsers from "@/components/chat/OnlineUsers";

interface User {
  id: string;
  username: string;
  display_name: string;
}

interface Room {
  id: string;
  name: string;
  description: string | null;
}

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001";

export default function ChatRoom() {
  const params = useParams();
  const router = useRouter();
  const roomId = params.roomId as string;
  
  const [user, setUser] = useState<User | null>(null);
  const [room, setRoom] = useState<Room | null>(null);
  const [showOnlineUsers, setShowOnlineUsers] = useState(true);

  // Load user from localStorage
  useEffect(() => {
    const savedUser = localStorage.getItem("chat_user");
    if (savedUser) {
      setUser(JSON.parse(savedUser));
    } else {
      router.push("/");
    }
  }, [router]);

  // Fetch room details and mark as read
  useEffect(() => {
    if (roomId) {
      fetchRoom();
    }
  }, [roomId]);

  // Mark messages as read when entering room
  useEffect(() => {
    if (roomId && user) {
      markAsRead();
    }
  }, [roomId, user]);

  const fetchRoom = async () => {
    try {
      const res = await fetch(`${API_URL}/api/rooms/${roomId}`);
      if (res.ok) {
        const data = await res.json();
        setRoom(data);
      }
    } catch (error) {
      console.error("Failed to fetch room:", error);
    }
  };

  const markAsRead = async () => {
    if (!user) return;
    try {
      await fetch(`${API_URL}/api/rooms/${roomId}/read?userId=${user.id}`, {
        method: "POST",
      });
    } catch (error) {
      console.error("Failed to mark as read:", error);
    }
  };

  // WebSocket connection
  const {
    messages,
    onlineUsers,
    typingUsers,
    isConnected,
    sendMessage,
    sendTyping,
    sendStopTyping,
  } = useWebSocket(
    roomId,
    user?.id || "",
    user?.username || "",
    user?.display_name || ""
  );

  // Mark as read when new messages arrive
  useEffect(() => {
    if (messages.length > 0 && user) {
      markAsRead();
    }
  }, [messages.length]);

  const handleSendMessage = useCallback(
    (content: string) => {
      if (content.trim() && user) {
        sendMessage(content);
      }
    },
    [sendMessage, user]
  );

  if (!user) {
    return (
      <div className="flex items-center justify-center h-[calc(100vh-200px)]">
        <div className="text-center">
          <div className="animate-spin text-4xl mb-4">🔄</div>
          <p className="text-[var(--color-earth-600)]">กำลังโหลด...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-4 h-[calc(100vh-140px)]">
      <div className="card-isan h-full flex flex-col">
        {/* Room Header */}
        <div className="card-isan-header flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button
              onClick={() => router.push("/")}
              className="hover:bg-white/20 p-2 rounded-lg transition"
            >
              ← กลับ
            </button>
            <div>
              <h2 className="font-bold">{room?.name || "ห้องแชท"}</h2>
              <p className="text-xs text-[var(--color-gold-100)]">
                {room?.description || ""}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            {/* Connection status */}
            <div className="flex items-center gap-2 text-sm">
              <span
                className={`w-2 h-2 rounded-full ${
                  isConnected ? "bg-[var(--color-paddy-400)]" : "bg-red-400"
                }`}
              ></span>
              {isConnected ? "เชื่อมต่อแล้ว" : "กำลังเชื่อมต่อ..."}
            </div>
            {/* Toggle online users */}
            <button
              onClick={() => setShowOnlineUsers(!showOnlineUsers)}
              className="hover:bg-white/20 p-2 rounded-lg transition"
              title="แสดง/ซ่อนผู้ใช้ออนไลน์"
            >
              👥 {onlineUsers.length}
            </button>
          </div>
        </div>

        {/* Chat Content */}
        <div className="flex-1 flex overflow-hidden">
          {/* Messages Area */}
          <div className="flex-1 flex flex-col">
            <MessageList
              messages={messages}
              currentUserId={user.id}
              typingUsers={typingUsers.filter((u) => u.user_id !== user.id)}
            />
            <MessageInput
              onSendMessage={handleSendMessage}
              onTyping={sendTyping}
              onStopTyping={sendStopTyping}
              disabled={!isConnected}
            />
          </div>

          {/* Online Users Sidebar */}
          {showOnlineUsers && (
            <div className="w-64 border-l border-[var(--color-earth-200)] hidden lg:block">
              <OnlineUsers users={onlineUsers} currentUserId={user.id} />
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
