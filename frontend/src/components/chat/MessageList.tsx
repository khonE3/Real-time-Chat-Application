"use client";

import { useEffect, useRef } from "react";
import { Message, TypingUser } from "@/types";

interface MessageListProps {
  messages: Message[];
  currentUserId: string;
  typingUsers: TypingUser[];
}

export default function MessageList({
  messages,
  currentUserId,
  typingUsers,
}: MessageListProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom on new messages
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, typingUsers]);

  const formatTime = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString("th-TH", {
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const formatDate = (timestamp: string) => {
    const date = new Date(timestamp);
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(yesterday.getDate() - 1);

    if (date.toDateString() === today.toDateString()) {
      return "วันนี้";
    } else if (date.toDateString() === yesterday.toDateString()) {
      return "เมื่อวาน";
    } else {
      return date.toLocaleDateString("th-TH", {
        day: "numeric",
        month: "short",
        year: "numeric",
      });
    }
  };

  // Group messages by date
  const groupedMessages: { [date: string]: Message[] } = {};
  messages.forEach((msg) => {
    const dateKey = new Date(msg.created_at).toDateString();
    if (!groupedMessages[dateKey]) {
      groupedMessages[dateKey] = [];
    }
    groupedMessages[dateKey].push(msg);
  });

  return (
    <div className="flex-1 overflow-y-auto p-4 space-y-4">
      {messages.length === 0 ? (
        <div className="flex flex-col items-center justify-center h-full text-[var(--color-earth-500)]">
          <span className="text-6xl mb-4">💬</span>
          <p className="text-lg font-medium">ยังไม่มีข้อความ</p>
          <p className="text-sm">เริ่มพูดคุยกันเลย!</p>
        </div>
      ) : (
        Object.entries(groupedMessages).map(([dateKey, dateMessages]) => (
          <div key={dateKey}>
            {/* Date separator */}
            <div className="flex items-center justify-center my-4">
              <div className="bg-[var(--color-earth-200)] text-[var(--color-earth-600)] text-xs px-3 py-1 rounded-full">
                {formatDate(dateMessages[0].created_at)}
              </div>
            </div>

            {/* Messages */}
            {dateMessages.map((message, index) => {
              const isSent = message.user_id === currentUserId;
              const showAvatar =
                !isSent &&
                (index === 0 ||
                  dateMessages[index - 1].user_id !== message.user_id);

              return (
                <div
                  key={message.id}
                  className={`flex items-end gap-2 mb-2 ${
                    isSent ? "justify-end" : "justify-start"
                  }`}
                >
                  {/* Avatar for received messages */}
                  {!isSent && (
                    <div className="w-8 h-8 flex-shrink-0">
                      {showAvatar && (
                        <div className="w-8 h-8 bg-gradient-to-br from-[var(--color-silk-400)] to-[var(--color-silk-600)] rounded-full flex items-center justify-center text-white text-sm font-bold">
                          {(message.display_name || message.username || "?")
                            .charAt(0)
                            .toUpperCase()}
                        </div>
                      )}
                    </div>
                  )}

                  <div
                    className={`max-w-[70%] ${isSent ? "items-end" : "items-start"}`}
                  >
                    {/* Username for received messages */}
                    {!isSent && showAvatar && (
                      <p className="text-xs text-[var(--color-earth-500)] mb-1 ml-1">
                        {message.display_name || message.username}
                      </p>
                    )}

                    {/* Message bubble */}
                    <div className={`message-bubble ${isSent ? "sent" : "received"}`}>
                      <p className="break-words">{message.content}</p>
                      <p
                        className={`text-xs mt-1 ${
                          isSent
                            ? "text-white/70"
                            : "text-[var(--color-earth-400)]"
                        }`}
                      >
                        {formatTime(message.created_at)}
                      </p>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        ))
      )}

      {/* Typing indicators */}
      {typingUsers.length > 0 && (
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 bg-[var(--color-earth-200)] rounded-full flex items-center justify-center">
            <div className="typing-indicator">
              <span></span>
              <span></span>
              <span></span>
            </div>
          </div>
          <p className="text-sm text-[var(--color-earth-500)]">
            {typingUsers.map((u) => u.display_name || u.username).join(", ")}{" "}
            กำลังพิมพ์...
          </p>
        </div>
      )}

      <div ref={messagesEndRef} />
    </div>
  );
}
