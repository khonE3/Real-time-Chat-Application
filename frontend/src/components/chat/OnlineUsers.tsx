"use client";

import { OnlineUser } from "@/types";

interface OnlineUsersProps {
  users: OnlineUser[];
  currentUserId: string;
}

export default function OnlineUsers({ users, currentUserId }: OnlineUsersProps) {
  return (
    <div className="h-full flex flex-col bg-[var(--color-earth-50)]">
      {/* Header */}
      <div className="p-3 border-b border-[var(--color-earth-200)]">
        <h3 className="font-semibold text-[var(--color-earth-700)] flex items-center gap-2">
          <span className="online-indicator"></span>
          ออนไลน์ ({users.length})
        </h3>
      </div>

      {/* User list */}
      <div className="flex-1 overflow-y-auto p-2">
        {users.length === 0 ? (
          <div className="text-center py-4 text-[var(--color-earth-500)]">
            <p className="text-2xl mb-2">👤</p>
            <p className="text-sm">ไม่มีผู้ใช้ออนไลน์</p>
          </div>
        ) : (
          <div className="space-y-1">
            {users.map((user) => {
              const isCurrentUser = user.user_id === currentUserId;

              return (
                <div
                  key={user.user_id}
                  className={`flex items-center gap-2 p-2 rounded-lg transition ${
                    isCurrentUser
                      ? "bg-[var(--color-gold-100)]"
                      : "hover:bg-[var(--color-earth-100)]"
                  }`}
                >
                  {/* Avatar */}
                  <div className="relative">
                    <div
                      className={`w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-bold ${
                        isCurrentUser
                          ? "bg-gradient-to-br from-[var(--color-gold-400)] to-[var(--color-gold-600)]"
                          : "bg-gradient-to-br from-[var(--color-silk-400)] to-[var(--color-silk-600)]"
                      }`}
                    >
                      {(user.display_name || user.username || "?")
                        .charAt(0)
                        .toUpperCase()}
                    </div>
                    {/* Online indicator */}
                    <span className="absolute -bottom-0.5 -right-0.5 w-3 h-3 bg-[var(--color-paddy-500)] rounded-full border-2 border-[var(--color-earth-50)]"></span>
                  </div>

                  {/* User info */}
                  <div className="flex-1 min-w-0">
                    <p
                      className={`text-sm font-medium truncate ${
                        isCurrentUser
                          ? "text-[var(--color-gold-700)]"
                          : "text-[var(--color-earth-700)]"
                      }`}
                    >
                      {user.display_name || user.username}
                      {isCurrentUser && (
                        <span className="text-xs ml-1">(คุณ)</span>
                      )}
                    </p>
                    <p className="text-xs text-[var(--color-earth-500)] truncate">
                      @{user.username}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Footer */}
      <div className="p-3 border-t border-[var(--color-earth-200)] bg-white">
        <p className="text-xs text-[var(--color-earth-500)] text-center">
          🟢 แสดงผู้ใช้ที่ออนไลน์ใน 5 นาทีล่าสุด
        </p>
      </div>
    </div>
  );
}
