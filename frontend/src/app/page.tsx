"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";

interface Room {
  id: string;
  name: string;
  description: string | null;
  is_private: boolean;
  member_count: number;
  online_count: number;
  unread_count: number;
}

interface User {
  id: string;
  username: string;
  display_name: string;
}

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001";

export default function Home() {
  const router = useRouter();
  const [rooms, setRooms] = useState<Room[]>([]);
  const [user, setUser] = useState<User | null>(null);
  const [username, setUsername] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [showLogin, setShowLogin] = useState(true);

  // Check for existing user in localStorage
  useEffect(() => {
    const savedUser = localStorage.getItem("chat_user");
    if (savedUser) {
      setUser(JSON.parse(savedUser));
      setShowLogin(false);
    }
  }, []);

  // Fetch rooms
  useEffect(() => {
    fetchRooms();
  }, []);

  // Refetch rooms when user changes to get unread counts
  useEffect(() => {
    if (user) {
      fetchRooms();
    }
  }, [user]);

  const fetchRooms = async () => {
    try {
      // Include userId to get unread counts if user is logged in
      const savedUser = localStorage.getItem("chat_user");
      let url = `${API_URL}/api/rooms`;
      if (savedUser) {
        const userData = JSON.parse(savedUser);
        url += `?userId=${userData.id}`;
      }
      const res = await fetch(url);
      if (res.ok) {
        const data = await res.json();
        setRooms(data);
      }
    } catch (error) {
      console.error("Failed to fetch rooms:", error);
    }
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!username.trim() || !displayName.trim()) return;

    setIsLoading(true);
    try {
      const res = await fetch(`${API_URL}/api/users`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: username.trim().toLowerCase(),
          display_name: displayName.trim(),
        }),
      });

      if (res.ok) {
        const userData = await res.json();
        setUser(userData);
        localStorage.setItem("chat_user", JSON.stringify(userData));
        setShowLogin(false);
      }
    } catch (error) {
      console.error("Login failed:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("chat_user");
    setUser(null);
    setShowLogin(true);
    setUsername("");
    setDisplayName("");
  };

  const joinRoom = (roomId: string) => {
    if (!user) {
      setShowLogin(true);
      return;
    }
    router.push(`/chat/${roomId}`);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Welcome Section */}
      <div className="text-center mb-8">
        <h2 className="text-3xl font-bold text-[var(--color-earth-800)] mb-2">
          🎋 ยินดีต้อนรับสู่แชทอีสาน
        </h2>
        <p className="text-[var(--color-earth-600)]">
          พูดคุยแลกเปลี่ยน เหมือนนั่งกินข้าวเหนียวริมโขง
        </p>
      </div>

      <div className="grid lg:grid-cols-3 gap-8">
        {/* Login/User Card */}
        <div className="lg:col-span-1">
          <div className="card-isan">
            <div className="card-isan-header">
              {user ? "👤 ข้อมูลผู้ใช้" : "🔐 เข้าสู่ระบบ"}
            </div>
            <div className="p-4">
              {showLogin && !user ? (
                <form onSubmit={handleLogin} className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-[var(--color-earth-700)] mb-1">
                      ชื่อผู้ใช้ (Username)
                    </label>
                    <input
                      type="text"
                      value={username}
                      onChange={(e) => setUsername(e.target.value)}
                      placeholder="somchai"
                      className="input-isan"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-[var(--color-earth-700)] mb-1">
                      ชื่อที่แสดง (Display Name)
                    </label>
                    <input
                      type="text"
                      value={displayName}
                      onChange={(e) => setDisplayName(e.target.value)}
                      placeholder="สมชาย"
                      className="input-isan"
                      required
                    />
                  </div>
                  <button
                    type="submit"
                    disabled={isLoading}
                    className="btn-gold w-full"
                  >
                    {isLoading ? "กำลังเข้าสู่ระบบ..." : "เข้าสู่ระบบ"}
                  </button>
                </form>
              ) : user ? (
                <div className="space-y-4">
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 bg-gradient-to-br from-[var(--color-gold-400)] to-[var(--color-gold-600)] rounded-full flex items-center justify-center text-white text-xl font-bold">
                      {user.display_name.charAt(0).toUpperCase()}
                    </div>
                    <div>
                      <p className="font-semibold text-[var(--color-earth-800)]">
                        {user.display_name}
                      </p>
                      <p className="text-sm text-[var(--color-earth-500)]">
                        @{user.username}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-[var(--color-paddy-600)]">
                    <span className="online-indicator"></span>
                    ออนไลน์
                  </div>
                  <button onClick={handleLogout} className="btn-silk w-full">
                    ออกจากระบบ
                  </button>
                </div>
              ) : null}
            </div>
          </div>

          {/* Stats Card */}
          <div className="card-isan mt-4">
            <div className="p-4">
              <h3 className="font-semibold text-[var(--color-earth-700)] mb-3">
                📊 สถิติ
              </h3>
              <div className="grid grid-cols-2 gap-4 text-center">
                <div className="bg-[var(--color-gold-50)] rounded-lg p-3">
                  <p className="text-2xl font-bold text-[var(--color-gold-600)]">
                    {rooms.length}
                  </p>
                  <p className="text-xs text-[var(--color-earth-600)]">ห้องแชท</p>
                </div>
                <div className="bg-[var(--color-paddy-50)] rounded-lg p-3">
                  <p className="text-2xl font-bold text-[var(--color-paddy-600)]">
                    {rooms.reduce((acc, r) => acc + (r.online_count || 0), 0)}
                  </p>
                  <p className="text-xs text-[var(--color-earth-600)]">ออนไลน์</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Rooms List */}
        <div className="lg:col-span-2">
          <div className="card-isan">
            <div className="card-isan-header flex items-center justify-between">
              <span>🏠 ห้องแชท</span>
              <button
                onClick={fetchRooms}
                className="text-sm bg-white/20 px-3 py-1 rounded-full hover:bg-white/30 transition"
              >
                🔄 รีเฟรช
              </button>
            </div>
            <div className="p-4">
              {rooms.length === 0 ? (
                <div className="text-center py-8 text-[var(--color-earth-500)]">
                  <p className="text-4xl mb-2">🏯</p>
                  <p>ยังไม่มีห้องแชท กำลังโหลด...</p>
                </div>
              ) : (
                <div className="space-y-3">
                  {rooms.map((room) => (
                    <div
                      key={room.id}
                      className="flex items-center justify-between p-4 bg-[var(--color-earth-50)] rounded-lg hover:bg-[var(--color-gold-50)] transition cursor-pointer border border-transparent hover:border-[var(--color-gold-300)]"
                      onClick={() => joinRoom(room.id)}
                    >
                      <div className="flex items-center gap-3">
                        <div className="relative">
                          <div className="w-10 h-10 bg-gradient-to-br from-[var(--color-gold-400)] to-[var(--color-gold-600)] rounded-lg flex items-center justify-center text-white">
                            {room.name.includes("🏯")
                              ? "🏯"
                              : room.name.includes("🛒")
                              ? "🛒"
                              : room.name.includes("☕")
                              ? "☕"
                              : "💬"}
                          </div>
                          {room.unread_count > 0 && (
                            <span className="absolute -top-2 -right-2 min-w-5 h-5 bg-[var(--color-silk-500)] text-white text-xs font-bold rounded-full flex items-center justify-center px-1">
                              {room.unread_count > 99 ? "99+" : room.unread_count}
                            </span>
                          )}
                        </div>
                        <div>
                          <h4 className="font-semibold text-[var(--color-earth-800)]">
                            {room.name}
                          </h4>
                          {room.description && (
                            <p className="text-sm text-[var(--color-earth-500)]">
                              {room.description}
                            </p>
                          )}
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="flex items-center gap-1 text-sm text-[var(--color-paddy-600)]">
                          <span className="online-indicator"></span>
                          {room.online_count || 0}
                        </div>
                        <p className="text-xs text-[var(--color-earth-400)]">
                          {room.member_count || 0} สมาชิก
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Feature Highlights */}
          <div className="grid md:grid-cols-3 gap-4 mt-4">
            <div className="card-isan p-4 text-center">
              <span className="text-3xl">⚡</span>
              <h4 className="font-semibold text-[var(--color-earth-700)] mt-2">
                เรียลไทม์
              </h4>
              <p className="text-sm text-[var(--color-earth-500)]">
                ส่งข้อความทันที ไม่ต้องรอ
              </p>
            </div>
            <div className="card-isan p-4 text-center">
              <span className="text-3xl">👥</span>
              <h4 className="font-semibold text-[var(--color-earth-700)] mt-2">
                ผู้ใช้ออนไลน์
              </h4>
              <p className="text-sm text-[var(--color-earth-500)]">
                เห็นว่าใครออนไลน์อยู่
              </p>
            </div>
            <div className="card-isan p-4 text-center">
              <span className="text-3xl">✍️</span>
              <h4 className="font-semibold text-[var(--color-earth-700)] mt-2">
                กำลังพิมพ์
              </h4>
              <p className="text-sm text-[var(--color-earth-500)]">
                รู้ว่าใครกำลังพิมพ์ข้อความ
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
