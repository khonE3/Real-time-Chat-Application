"use client";

import { useState, useRef, useCallback, useEffect } from "react";

interface MessageInputProps {
  onSendMessage: (content: string) => void;
  onTyping: () => void;
  onStopTyping: () => void;
  disabled?: boolean;
}

export default function MessageInput({
  onSendMessage,
  onTyping,
  onStopTyping,
  disabled = false,
}: MessageInputProps) {
  const [message, setMessage] = useState("");
  const [isTyping, setIsTyping] = useState(false);
  const typingTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const inputRef = useRef<HTMLTextAreaElement>(null);

  // Handle typing indicator
  const handleTyping = useCallback(() => {
    if (!isTyping) {
      setIsTyping(true);
      onTyping();
    }

    // Clear existing timeout
    if (typingTimeoutRef.current) {
      clearTimeout(typingTimeoutRef.current);
    }

    // Set new timeout to stop typing
    typingTimeoutRef.current = setTimeout(() => {
      setIsTyping(false);
      onStopTyping();
    }, 2000);
  }, [isTyping, onTyping, onStopTyping]);

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (typingTimeoutRef.current) {
        clearTimeout(typingTimeoutRef.current);
      }
    };
  }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (message.trim() && !disabled) {
      onSendMessage(message.trim());
      setMessage("");
      setIsTyping(false);
      onStopTyping();

      if (typingTimeoutRef.current) {
        clearTimeout(typingTimeoutRef.current);
      }
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setMessage(e.target.value);
    if (e.target.value.trim()) {
      handleTyping();
    }
  };

  // Auto-resize textarea
  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.style.height = "auto";
      inputRef.current.style.height = `${Math.min(
        inputRef.current.scrollHeight,
        120
      )}px`;
    }
  }, [message]);

  return (
    <form
      onSubmit={handleSubmit}
      className="border-t border-[var(--color-earth-200)] p-4 bg-white"
    >
      <div className="flex items-end gap-3">
        {/* Emoji button (placeholder) */}
        <button
          type="button"
          className="text-2xl hover:scale-110 transition"
          title="อิโมจิ (เร็วๆ นี้)"
        >
          😊
        </button>

        {/* Input area */}
        <div className="flex-1 relative">
          <textarea
            ref={inputRef}
            value={message}
            onChange={handleChange}
            onKeyDown={handleKeyDown}
            placeholder={disabled ? "กำลังเชื่อมต่อ..." : "พิมพ์ข้อความ..."}
            disabled={disabled}
            rows={1}
            className="input-isan resize-none pr-12"
            style={{ minHeight: "44px" }}
          />
          <span className="absolute right-3 bottom-3 text-xs text-[var(--color-earth-400)]">
            {message.length}/4000
          </span>
        </div>

        {/* Send button */}
        <button
          type="submit"
          disabled={!message.trim() || disabled}
          className={`btn-gold px-4 py-3 text-lg transition ${
            !message.trim() || disabled
              ? "opacity-50 cursor-not-allowed"
              : "hover:scale-105"
          }`}
        >
          📤
        </button>
      </div>

      {/* Hints */}
      <p className="text-xs text-[var(--color-earth-400)] mt-2">
        กด Enter เพื่อส่ง • Shift+Enter เพื่อขึ้นบรรทัดใหม่
      </p>
    </form>
  );
}
