'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

interface User {
    id: string;
    username: string;
    display_name: string;
    avatar_url?: string;
}

interface DMRoom {
    id: string;
    name: string;
    is_dm: boolean;
    other_user?: User;
    created_at: string;
}

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3001';

export default function DMPage() {
    const [user, setUser] = useState<User | null>(null);
    const [dmRooms, setDmRooms] = useState<DMRoom[]>([]);
    const [users, setUsers] = useState<User[]>([]);
    const [searchQuery, setSearchQuery] = useState('');
    const [isLoading, setIsLoading] = useState(true);
    const router = useRouter();

    useEffect(() => {
        const storedUser = localStorage.getItem('user');
        if (storedUser) {
            const parsed = JSON.parse(storedUser);
            setUser(parsed);
            fetchDMs(parsed.id);
        } else {
            router.push('/');
        }
    }, [router]);

    const fetchDMs = async (userId: string) => {
        setIsLoading(true);
        try {
            const res = await fetch(`${API_BASE}/api/dm?userId=${userId}`);
            if (res.ok) {
                const data = await res.json();
                setDmRooms(data || []);
            }
        } catch (error) {
            console.error('Failed to fetch DMs:', error);
        } finally {
            setIsLoading(false);
        }
    };

    const searchUsers = async (query: string) => {
        if (query.length < 2) {
            setUsers([]);
            return;
        }

        try {
            const res = await fetch(`${API_BASE}/api/search/users?q=${encodeURIComponent(query)}`);
            if (res.ok) {
                const data = await res.json();
                setUsers(data.filter((u: User) => u.id !== user?.id) || []);
            }
        } catch (error) {
            console.error('Failed to search users:', error);
        }
    };

    const startDM = async (targetUserId: string) => {
        if (!user) return;

        try {
            const res = await fetch(`${API_BASE}/api/dm/start?userId=${user.id}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ target_user_id: targetUserId }),
            });

            if (res.ok) {
                const data = await res.json();
                router.push(`/chat/${data.room.id}`);
            }
        } catch (error) {
            console.error('Failed to start DM:', error);
        }
    };

    const formatTime = (dateStr: string) => {
        const date = new Date(dateStr);
        const now = new Date();
        const diff = now.getTime() - date.getTime();

        if (diff < 86400000) return date.toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' });
        if (diff < 604800000) return date.toLocaleDateString('th-TH', { weekday: 'short' });
        return date.toLocaleDateString('th-TH', { month: 'short', day: 'numeric' });
    };

    if (!user) return null;

    return (
        <div className="dm-page">
            <header className="header">
                <Link href="/" className="back-btn">
                    ← กลับ
                </Link>
                <h1>ข้อความส่วนตัว</h1>
            </header>

            <div className="search-section">
                <input
                    type="text"
                    placeholder="ค้นหาผู้ใช้..."
                    value={searchQuery}
                    onChange={(e) => {
                        setSearchQuery(e.target.value);
                        searchUsers(e.target.value);
                    }}
                    className="search-input"
                />

                {users.length > 0 && (
                    <div className="user-list">
                        {users.map((u) => (
                            <button
                                key={u.id}
                                onClick={() => startDM(u.id)}
                                className="user-item"
                            >
                                <span className="avatar">
                                    {u.avatar_url ? (
                                        <img src={u.avatar_url} alt="" />
                                    ) : (
                                        u.display_name[0].toUpperCase()
                                    )}
                                </span>
                                <div className="user-info">
                                    <span className="name">{u.display_name}</span>
                                    <span className="username">@{u.username}</span>
                                </div>
                            </button>
                        ))}
                    </div>
                )}
            </div>

            <div className="dm-list">
                <h2>การสนทนาของคุณ</h2>

                {isLoading ? (
                    <div className="loading">กำลังโหลด...</div>
                ) : dmRooms.length === 0 ? (
                    <div className="empty">
                        <p>ยังไม่มีข้อความส่วนตัว</p>
                        <p>ค้นหาผู้ใช้เพื่อเริ่มสนทนา</p>
                    </div>
                ) : (
                    dmRooms.map((dm) => (
                        <Link
                            key={dm.id}
                            href={`/chat/${dm.id}`}
                            className="dm-item"
                        >
                            <span className="avatar">
                                {dm.other_user?.avatar_url ? (
                                    <img src={dm.other_user.avatar_url} alt="" />
                                ) : (
                                    dm.other_user?.display_name?.[0]?.toUpperCase() || '?'
                                )}
                            </span>
                            <div className="dm-info">
                                <span className="name">
                                    {dm.other_user?.display_name || 'Unknown User'}
                                </span>
                                <span className="username">
                                    @{dm.other_user?.username || 'unknown'}
                                </span>
                            </div>
                            <span className="time">{formatTime(dm.created_at)}</span>
                        </Link>
                    ))
                )}
            </div>

            <style jsx>{`
        .dm-page {
          min-height: 100vh;
          background: var(--color-earth-50, #faf8f5);
        }
        .header {
          display: flex;
          align-items: center;
          gap: 16px;
          padding: 16px 24px;
          background: white;
          border-bottom: 1px solid var(--color-earth-200, #e8dcc8);
        }
        .back-btn {
          color: var(--color-earth-600, #8b7355);
          text-decoration: none;
          font-size: 0.9rem;
        }
        .back-btn:hover {
          color: var(--color-gold-600, #c99a20);
        }
        .header h1 {
          margin: 0;
          font-size: 1.25rem;
          color: var(--color-earth-800, #4a3f2f);
        }
        .search-section {
          padding: 16px 24px;
          background: white;
          border-bottom: 1px solid var(--color-earth-200, #e8dcc8);
        }
        .search-input {
          width: 100%;
          padding: 12px 16px;
          border: 1px solid var(--color-earth-200, #e8dcc8);
          border-radius: 24px;
          font-size: 0.95rem;
          outline: none;
          transition: border-color 0.2s;
        }
        .search-input:focus {
          border-color: var(--color-gold-400, #f5c542);
        }
        .user-list {
          margin-top: 12px;
          background: white;
          border-radius: 12px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
          overflow: hidden;
        }
        .user-item {
          display: flex;
          align-items: center;
          gap: 12px;
          width: 100%;
          padding: 12px 16px;
          border: none;
          background: transparent;
          cursor: pointer;
          text-align: left;
          transition: background 0.2s;
        }
        .user-item:hover {
          background: var(--color-gold-50, #fffef5);
        }
        .avatar {
          width: 48px;
          height: 48px;
          border-radius: 50%;
          background: linear-gradient(135deg, var(--color-gold-400, #f5c542), var(--color-gold-500, #e5b32a));
          display: flex;
          align-items: center;
          justify-content: center;
          font-weight: 600;
          font-size: 1.25rem;
          color: var(--color-earth-900, #2d2416);
          overflow: hidden;
          flex-shrink: 0;
        }
        .avatar img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
        .user-info,
        .dm-info {
          flex: 1;
          min-width: 0;
        }
        .name {
          display: block;
          font-weight: 600;
          color: var(--color-earth-800, #4a3f2f);
        }
        .username {
          display: block;
          font-size: 0.85rem;
          color: var(--color-earth-500, #a08060);
        }
        .dm-list {
          padding: 24px;
        }
        .dm-list h2 {
          margin: 0 0 16px;
          font-size: 0.9rem;
          font-weight: 600;
          color: var(--color-earth-600, #8b7355);
          text-transform: uppercase;
        }
        .loading,
        .empty {
          text-align: center;
          padding: 48px 24px;
          color: var(--color-earth-500, #a08060);
        }
        .dm-item {
          display: flex;
          align-items: center;
          gap: 12px;
          padding: 16px;
          background: white;
          border-radius: 12px;
          margin-bottom: 8px;
          text-decoration: none;
          transition: all 0.2s;
        }
        .dm-item:hover {
          transform: translateX(4px);
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }
        .time {
          font-size: 0.8rem;
          color: var(--color-earth-400, #b0a090);
        }
      `}</style>
        </div>
    );
}
