"use client";

import { useState, useEffect, useRef, useCallback } from "react";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:3001";

// Types for global WebSocket messages
interface RoomStatsPayload {
    room_id: string;
    online_count: number;
    has_new_msg?: boolean;
    sender_id?: string;
}

interface GlobalPresencePayload {
    total_online: number;
}

interface Room {
    id: string;
    name: string;
    description: string | null;
    is_private: boolean;
    member_count: number;
    online_count: number;
    unread_count: number;
}

interface GlobalWSMessage {
    type: "room_stats" | "room_created" | "room_deleted" | "global_presence" | "rooms_init";
    payload: unknown;
}

interface UseGlobalWebSocketReturn {
    isConnected: boolean;
    roomUpdates: Map<string, { online_count: number; has_new_msg?: boolean }>;
    newRooms: Room[];
    totalOnline: number;
    onRoomUpdate: (callback: (roomId: string, stats: RoomStatsPayload) => void) => void;
}

export function useGlobalWebSocket(userId: string): UseGlobalWebSocketReturn {
    const [isConnected, setIsConnected] = useState(false);
    const [roomUpdates, setRoomUpdates] = useState<Map<string, { online_count: number; has_new_msg?: boolean }>>(new Map());
    const [newRooms, setNewRooms] = useState<Room[]>([]);
    const [totalOnline, setTotalOnline] = useState(0);

    const wsRef = useRef<WebSocket | null>(null);
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const reconnectAttempts = useRef(0);
    const callbackRef = useRef<((roomId: string, stats: RoomStatsPayload) => void) | null>(null);
    const maxReconnectAttempts = 5;

    const connect = useCallback(() => {
        if (!userId || userId === "") {
            return;
        }

        // Close existing connection
        if (wsRef.current) {
            wsRef.current.onclose = null;
            wsRef.current.close();
            wsRef.current = null;
        }

        if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current);
            reconnectTimeoutRef.current = null;
        }

        const wsUrl = `${WS_URL}/ws/global?userId=${userId}`;
        console.log("🌐 Connecting to Global WebSocket:", wsUrl);

        try {
            const ws = new WebSocket(wsUrl);
            wsRef.current = ws;

            ws.onopen = () => {
                console.log("✅ Global WebSocket connected");
                setIsConnected(true);
                reconnectAttempts.current = 0;
            };

            ws.onclose = (event) => {
                console.log("❌ Global WebSocket disconnected:", event.code);
                setIsConnected(false);

                if (wsRef.current === ws && reconnectAttempts.current < maxReconnectAttempts) {
                    reconnectAttempts.current++;
                    const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.current), 30000);
                    console.log(`🔄 Global WS reconnecting in ${delay}ms (attempt ${reconnectAttempts.current})`);

                    reconnectTimeoutRef.current = setTimeout(() => {
                        connect();
                    }, delay);
                }
            };

            ws.onerror = (event) => {
                console.error("⚠️ Global WebSocket error:", event);
            };

            ws.onmessage = (event) => {
                try {
                    const data: GlobalWSMessage = JSON.parse(event.data);
                    handleMessage(data);
                } catch (e) {
                    console.error("Failed to parse global message:", e);
                }
            };
        } catch (e) {
            console.error("Failed to create Global WebSocket:", e);
        }
    }, [userId]);

    const handleMessage = useCallback((data: GlobalWSMessage) => {
        switch (data.type) {
            case "room_stats":
                const stats = data.payload as RoomStatsPayload;
                setRoomUpdates((prev) => {
                    const newMap = new Map(prev);
                    newMap.set(stats.room_id, {
                        online_count: stats.online_count,
                        has_new_msg: stats.has_new_msg,
                    });
                    return newMap;
                });

                // Call external callback if registered
                if (callbackRef.current) {
                    callbackRef.current(stats.room_id, stats);
                }
                break;

            case "room_created":
                const newRoom = data.payload as Room;
                setNewRooms((prev) => [...prev, newRoom]);
                break;

            case "global_presence":
                const presence = data.payload as GlobalPresencePayload;
                setTotalOnline(presence.total_online);
                break;

            case "rooms_init":
                // Initial room list - used for syncing
                console.log("📋 Received initial room list via WebSocket");
                break;

            default:
                console.log("Unknown global message type:", data.type);
        }
    }, []);

    const onRoomUpdate = useCallback((callback: (roomId: string, stats: RoomStatsPayload) => void) => {
        callbackRef.current = callback;
    }, []);

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
        isConnected,
        roomUpdates,
        newRooms,
        totalOnline,
        onRoomUpdate,
    };
}
