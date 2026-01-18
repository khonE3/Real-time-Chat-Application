'use client';

import { useState, useEffect, useRef } from 'react';

interface Notification {
    id: string;
    type: string;
    title: string;
    body: string;
    is_read: boolean;
    created_at: string;
    from_username?: string;
    from_display_name?: string;
}

interface NotificationBellProps {
    userId: string;
}

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3001';

export default function NotificationBell({ userId }: NotificationBellProps) {
    const [isOpen, setIsOpen] = useState(false);
    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [unreadCount, setUnreadCount] = useState(0);
    const [isLoading, setIsLoading] = useState(false);
    const dropdownRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        fetchUnreadCount();
        const interval = setInterval(fetchUnreadCount, 30000); // Poll every 30s
        return () => clearInterval(interval);
    }, [userId]);

    useEffect(() => {
        const handleClickOutside = (e: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
                setIsOpen(false);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const fetchUnreadCount = async () => {
        try {
            const res = await fetch(`${API_BASE}/api/notifications/count?userId=${userId}`);
            if (res.ok) {
                const data = await res.json();
                setUnreadCount(data.unread_count || 0);
            }
        } catch (error) {
            console.error('Failed to fetch unread count:', error);
        }
    };

    const fetchNotifications = async () => {
        setIsLoading(true);
        try {
            const res = await fetch(`${API_BASE}/api/notifications?userId=${userId}&limit=10`);
            if (res.ok) {
                const data = await res.json();
                setNotifications(data.notifications || []);
                setUnreadCount(data.unread_count || 0);
            }
        } catch (error) {
            console.error('Failed to fetch notifications:', error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleOpen = () => {
        setIsOpen(!isOpen);
        if (!isOpen) {
            fetchNotifications();
        }
    };

    const markAsRead = async (id: string) => {
        try {
            await fetch(`${API_BASE}/api/notifications/${id}/read?userId=${userId}`, {
                method: 'POST',
            });
            setNotifications((prev) =>
                prev.map((n) => (n.id === id ? { ...n, is_read: true } : n))
            );
            setUnreadCount((prev) => Math.max(0, prev - 1));
        } catch (error) {
            console.error('Failed to mark as read:', error);
        }
    };

    const markAllAsRead = async () => {
        try {
            await fetch(`${API_BASE}/api/notifications/read-all?userId=${userId}`, {
                method: 'POST',
            });
            setNotifications((prev) => prev.map((n) => ({ ...n, is_read: true })));
            setUnreadCount(0);
        } catch (error) {
            console.error('Failed to mark all as read:', error);
        }
    };

    const formatTime = (dateStr: string) => {
        const date = new Date(dateStr);
        const now = new Date();
        const diff = now.getTime() - date.getTime();

        if (diff < 60000) return 'เมื่อกี้';
        if (diff < 3600000) return `${Math.floor(diff / 60000)} นาทีที่แล้ว`;
        if (diff < 86400000) return `${Math.floor(diff / 3600000)} ชั่วโมงที่แล้ว`;
        return `${Math.floor(diff / 86400000)} วันที่แล้ว`;
    };

    const getIcon = (type: string) => {
        switch (type) {
            case 'message': return '💬';
            case 'mention': return '📢';
            case 'reaction': return '❤️';
            case 'dm': return '✉️';
            default: return '🔔';
        }
    };

    return (
        <div ref={dropdownRef} className="notification-bell">
            <button onClick={handleOpen} className="bell-btn">
                <span className="icon">🔔</span>
                {unreadCount > 0 && (
                    <span className="badge">{unreadCount > 9 ? '9+' : unreadCount}</span>
                )}
            </button>

            {isOpen && (
                <div className="dropdown">
                    <div className="header">
                        <h3>การแจ้งเตือน</h3>
                        {unreadCount > 0 && (
                            <button onClick={markAllAsRead} className="mark-all">
                                อ่านทั้งหมด
                            </button>
                        )}
                    </div>

                    <div className="list">
                        {isLoading ? (
                            <div className="loading">กำลังโหลด...</div>
                        ) : notifications.length === 0 ? (
                            <div className="empty">ไม่มีการแจ้งเตือน</div>
                        ) : (
                            notifications.map((notification) => (
                                <div
                                    key={notification.id}
                                    className={`item ${!notification.is_read ? 'unread' : ''}`}
                                    onClick={() => !notification.is_read && markAsRead(notification.id)}
                                >
                                    <span className="type-icon">{getIcon(notification.type)}</span>
                                    <div className="content">
                                        <p className="title">{notification.title}</p>
                                        {notification.body && (
                                            <p className="body">{notification.body}</p>
                                        )}
                                        <p className="time">{formatTime(notification.created_at)}</p>
                                    </div>
                                    {!notification.is_read && <span className="dot" />}
                                </div>
                            ))
                        )}
                    </div>
                </div>
            )}

            <style jsx>{`
        .notification-bell {
          position: relative;
        }
        .bell-btn {
          position: relative;
          padding: 8px;
          background: transparent;
          border: none;
          cursor: pointer;
          font-size: 1.25rem;
          border-radius: 8px;
          transition: background 0.2s;
        }
        .bell-btn:hover {
          background: var(--color-earth-100, #f5f0e8);
        }
        .badge {
          position: absolute;
          top: 0;
          right: 0;
          background: var(--color-silk-500, #e63946);
          color: white;
          font-size: 0.65rem;
          font-weight: 600;
          padding: 2px 5px;
          border-radius: 10px;
          min-width: 16px;
          text-align: center;
        }
        .dropdown {
          position: absolute;
          top: 100%;
          right: 0;
          margin-top: 8px;
          background: white;
          border-radius: 12px;
          box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
          width: 320px;
          max-height: 400px;
          overflow: hidden;
          z-index: 1000;
        }
        .header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 12px 16px;
          border-bottom: 1px solid var(--color-earth-200, #e8dcc8);
        }
        .header h3 {
          margin: 0;
          font-size: 1rem;
          color: var(--color-earth-800, #4a3f2f);
        }
        .mark-all {
          font-size: 0.8rem;
          color: var(--color-gold-600, #c99a20);
          background: none;
          border: none;
          cursor: pointer;
        }
        .mark-all:hover {
          text-decoration: underline;
        }
        .list {
          max-height: 340px;
          overflow-y: auto;
        }
        .loading,
        .empty {
          padding: 24px;
          text-align: center;
          color: var(--color-earth-500, #a08060);
        }
        .item {
          display: flex;
          gap: 12px;
          padding: 12px 16px;
          cursor: pointer;
          transition: background 0.2s;
          position: relative;
        }
        .item:hover {
          background: var(--color-earth-50, #faf8f5);
        }
        .item.unread {
          background: var(--color-gold-50, #fffef5);
        }
        .type-icon {
          font-size: 1.25rem;
          flex-shrink: 0;
        }
        .content {
          flex: 1;
          min-width: 0;
        }
        .title {
          margin: 0 0 4px;
          font-size: 0.9rem;
          font-weight: 500;
          color: var(--color-earth-800, #4a3f2f);
        }
        .body {
          margin: 0 0 4px;
          font-size: 0.85rem;
          color: var(--color-earth-600, #8b7355);
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        .time {
          margin: 0;
          font-size: 0.75rem;
          color: var(--color-earth-400, #b0a090);
        }
        .dot {
          position: absolute;
          right: 16px;
          top: 50%;
          transform: translateY(-50%);
          width: 8px;
          height: 8px;
          background: var(--color-gold-400, #f5c542);
          border-radius: 50%;
        }
      `}</style>
        </div>
    );
}
