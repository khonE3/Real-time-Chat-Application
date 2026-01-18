'use client';

import { useState, useCallback, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';

interface SearchBarProps {
    roomId?: string;
    placeholder?: string;
}

interface SearchResult {
    users?: Array<{
        id: string;
        username: string;
        display_name: string;
        avatar_url?: string;
    }>;
    messages?: Array<{
        id: string;
        content: string;
        username: string;
        created_at: string;
    }>;
}

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3001';

export default function SearchBar({ roomId, placeholder = 'ค้นหา...' }: SearchBarProps) {
    const [query, setQuery] = useState('');
    const [isOpen, setIsOpen] = useState(false);
    const [results, setResults] = useState<SearchResult>({});
    const [isLoading, setIsLoading] = useState(false);
    const searchRef = useRef<HTMLDivElement>(null);
    const debounceRef = useRef<NodeJS.Timeout>();
    const router = useRouter();

    useEffect(() => {
        const handleClickOutside = (e: MouseEvent) => {
            if (searchRef.current && !searchRef.current.contains(e.target as Node)) {
                setIsOpen(false);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const search = useCallback(async (q: string) => {
        if (q.length < 2) {
            setResults({});
            return;
        }

        setIsLoading(true);
        try {
            const params = new URLSearchParams({ q });
            if (roomId) params.append('room', roomId);

            const res = await fetch(`${API_BASE}/api/search?${params}`);
            if (res.ok) {
                const data = await res.json();
                setResults(data);
                setIsOpen(true);
            }
        } catch (error) {
            console.error('Search failed:', error);
        } finally {
            setIsLoading(false);
        }
    }, [roomId]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        setQuery(value);

        // Debounce search
        if (debounceRef.current) {
            clearTimeout(debounceRef.current);
        }

        debounceRef.current = setTimeout(() => {
            search(value);
        }, 300);
    };

    const handleUserClick = (userId: string) => {
        router.push(`/dm/${userId}`);
        setIsOpen(false);
        setQuery('');
    };

    const handleMessageClick = (messageId: string) => {
        // Could scroll to message or navigate
        console.log('Navigate to message:', messageId);
        setIsOpen(false);
    };

    const formatTime = (dateStr: string) => {
        const date = new Date(dateStr);
        return date.toLocaleDateString('th-TH', {
            month: 'short',
            day: 'numeric',
        });
    };

    const hasResults = (results.users?.length || 0) + (results.messages?.length || 0) > 0;

    return (
        <div ref={searchRef} className="search-bar">
            <div className="input-wrapper">
                <span className="icon">🔍</span>
                <input
                    type="text"
                    value={query}
                    onChange={handleChange}
                    onFocus={() => hasResults && setIsOpen(true)}
                    placeholder={placeholder}
                    className="input"
                />
                {isLoading && <span className="loading">⏳</span>}
                {query && !isLoading && (
                    <button
                        onClick={() => {
                            setQuery('');
                            setResults({});
                        }}
                        className="clear"
                    >
                        ✕
                    </button>
                )}
            </div>

            {isOpen && hasResults && (
                <div className="dropdown">
                    {results.users && results.users.length > 0 && (
                        <div className="section">
                            <h4 className="section-title">ผู้ใช้</h4>
                            {results.users.map((user) => (
                                <button
                                    key={user.id}
                                    onClick={() => handleUserClick(user.id)}
                                    className="result-item"
                                >
                                    <span className="avatar">
                                        {user.avatar_url ? (
                                            <img src={user.avatar_url} alt="" />
                                        ) : (
                                            user.display_name[0].toUpperCase()
                                        )}
                                    </span>
                                    <div className="info">
                                        <span className="name">{user.display_name}</span>
                                        <span className="username">@{user.username}</span>
                                    </div>
                                </button>
                            ))}
                        </div>
                    )}

                    {results.messages && results.messages.length > 0 && (
                        <div className="section">
                            <h4 className="section-title">ข้อความ</h4>
                            {results.messages.map((msg) => (
                                <button
                                    key={msg.id}
                                    onClick={() => handleMessageClick(msg.id)}
                                    className="result-item"
                                >
                                    <span className="msg-icon">💬</span>
                                    <div className="info">
                                        <span className="content">{msg.content}</span>
                                        <span className="meta">
                                            {msg.username} • {formatTime(msg.created_at)}
                                        </span>
                                    </div>
                                </button>
                            ))}
                        </div>
                    )}
                </div>
            )}

            <style jsx>{`
        .search-bar {
          position: relative;
          width: 100%;
          max-width: 400px;
        }
        .input-wrapper {
          display: flex;
          align-items: center;
          gap: 8px;
          padding: 8px 12px;
          background: var(--color-earth-50, #faf8f5);
          border: 1px solid var(--color-earth-200, #e8dcc8);
          border-radius: 24px;
          transition: all 0.2s;
        }
        .input-wrapper:focus-within {
          border-color: var(--color-gold-400, #f5c542);
          background: white;
        }
        .icon {
          font-size: 1rem;
          opacity: 0.6;
        }
        .input {
          flex: 1;
          border: none;
          background: transparent;
          outline: none;
          font-size: 0.9rem;
          color: var(--color-earth-800, #4a3f2f);
        }
        .input::placeholder {
          color: var(--color-earth-400, #b0a090);
        }
        .loading {
          font-size: 0.9rem;
        }
        .clear {
          width: 18px;
          height: 18px;
          border: none;
          background: var(--color-earth-200, #e8dcc8);
          border-radius: 50%;
          cursor: pointer;
          font-size: 0.65rem;
          color: var(--color-earth-600, #8b7355);
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .clear:hover {
          background: var(--color-earth-300, #d4c4a8);
        }
        .dropdown {
          position: absolute;
          top: 100%;
          left: 0;
          right: 0;
          margin-top: 8px;
          background: white;
          border-radius: 12px;
          box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
          max-height: 400px;
          overflow-y: auto;
          z-index: 1000;
        }
        .section {
          padding: 8px 0;
        }
        .section:not(:last-child) {
          border-bottom: 1px solid var(--color-earth-200, #e8dcc8);
        }
        .section-title {
          margin: 0;
          padding: 6px 16px;
          font-size: 0.75rem;
          font-weight: 600;
          color: var(--color-earth-500, #a08060);
          text-transform: uppercase;
        }
        .result-item {
          display: flex;
          align-items: center;
          gap: 12px;
          width: 100%;
          padding: 10px 16px;
          border: none;
          background: transparent;
          cursor: pointer;
          text-align: left;
          transition: background 0.2s;
        }
        .result-item:hover {
          background: var(--color-gold-50, #fffef5);
        }
        .avatar {
          width: 36px;
          height: 36px;
          border-radius: 50%;
          background: linear-gradient(135deg, var(--color-gold-400, #f5c542), var(--color-gold-500, #e5b32a));
          display: flex;
          align-items: center;
          justify-content: center;
          font-weight: 600;
          color: var(--color-earth-900, #2d2416);
          overflow: hidden;
        }
        .avatar img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
        .msg-icon {
          font-size: 1.25rem;
        }
        .info {
          flex: 1;
          min-width: 0;
        }
        .name {
          display: block;
          font-weight: 500;
          color: var(--color-earth-800, #4a3f2f);
        }
        .username,
        .meta {
          display: block;
          font-size: 0.8rem;
          color: var(--color-earth-500, #a08060);
        }
        .content {
          display: block;
          font-size: 0.9rem;
          color: var(--color-earth-700, #6b5c48);
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      `}</style>
        </div>
    );
}
