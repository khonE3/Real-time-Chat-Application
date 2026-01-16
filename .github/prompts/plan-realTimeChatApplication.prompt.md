## Plan: Real-time Chat Application with Redis + Go Fiber + PostgreSQL

สร้างแอปพลิเคชันแชทเรียลไทม์ ใช้ Go Fiber จัดการ WebSocket, Redis 7.4 (Pub/Sub, Streams, Online Tracking) สำหรับ real-time features และ PostgreSQL สำหรับ persistent data พร้อม Frontend Next.js 16 ธีมอีสานหนองบัวลำภู

### Tech Stack Versions
| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.25.6 | Backend runtime |
| Go Fiber | v3.0.0-rc.3 | Web framework + WebSocket |
| go-redis | v9.17.2 | Redis client |
| pgx | v5.7.2 | PostgreSQL driver |
| Next.js | 16.1 | Frontend framework |
| React | 19.2 | UI library |
| TailwindCSS | 4.1 | Styling |
| Bun | 1.3.6 | Package manager |
| Redis | 7.4 | Real-time (Pub/Sub, Streams, Presence) |
| PostgreSQL | 17.2 | Persistent storage |

### Architecture Overview
```
┌─────────────────────────────────────────────────────────────────┐
│                      DATA LAYER                                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────────────┐    ┌──────────────────────┐           │
│  │      PostgreSQL      │    │        Redis         │           │
│  │   (Persistent Data)  │    │   (Real-time Data)   │           │
│  ├──────────────────────┤    ├──────────────────────┤           │
│  │ • Users              │    │ • Pub/Sub (live msg) │           │
│  │ • Rooms              │    │ • Streams (recent)   │           │
│  │ • Messages (history) │    │ • Sorted Sets        │           │
│  │ • Room members       │    │   (online users)     │           │
│  └──────────────────────┘    └──────────────────────┘           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Features Checklist

#### ✅ Core Features (MVP)
- [ ] Real-time messaging (WebSocket + Redis Pub/Sub)
- [ ] Message history (PostgreSQL + Redis Streams)
- [ ] Online user tracking (Redis Sorted Sets)
- [ ] Multiple chat rooms (public)
- [ ] Isan Nong Bua Lam Phu theme UI
- [ ] Simple Auth (nickname-based)
- [ ] Typing indicators (Redis Pub/Sub)
- [ ] Unread message count

#### 🔵 Nice to Have (Phase 2)
- [ ] JWT authentication
- [ ] File/Image upload (MinIO/S3)
- [ ] Edit/Delete messages
- [ ] Reply to message
- [ ] Emoji reactions
- [ ] Browser push notifications
- [ ] Message search (PostgreSQL full-text)
- [ ] User profiles & avatars
- [ ] Dark mode toggle

### Steps

1. **สร้างโครงสร้างโปรเจค** - แยก `backend/` (Go Fiber) และ `frontend/` (Next.js) พร้อม README.md, .gitignore และ docker-compose.yml สำหรับ PostgreSQL

2. **พัฒนา Backend Go Fiber** - สร้าง WebSocket hub ใน `backend/internal/websocket/hub.go`, Redis client ใน `backend/pkg/redis/client.go`, PostgreSQL connection ใน `backend/pkg/database/postgres.go`

3. **สร้าง Database Schema (PostgreSQL)** - ออกแบบ tables: `users`, `rooms`, `room_members`, `messages` พร้อม migrations ใน `backend/migrations/`

4. **สร้าง Redis Repository Layer** - Implement Pub/Sub (`pubsub_repo.go`), Streams สำหรับ recent messages (`stream_repo.go`), และ Sorted Sets สำหรับ online tracking (`presence_repo.go`)

5. **สร้าง PostgreSQL Repository Layer** - Implement `user_repo.go`, `room_repo.go`, `message_repo.go` สำหรับ persistent CRUD operations

6. **พัฒนา Frontend Next.js** - สร้าง WebSocket hook (`useWebSocket.ts`), Chat components (`ChatRoom.tsx`, `MessageList.tsx`, `MessageInput.tsx`, `OnlineUsers.tsx`)

7. **ออกแบบธีมอีสานหนองบัวลำภู** - กำหนด color palette (สีทองวัด `#D4A12A`, แดงผ้าไหม `#DC143C`, เขียวทุ่งนา `#228B22`, น้ำตาลดิน `#8B4513`)

8. **สร้าง Documentation** - เขียน README.md อธิบาย architecture, วิธี setup, database schema และ API endpoints

### Database Schema (PostgreSQL)

```sql
-- users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- rooms table
CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_private BOOLEAN DEFAULT FALSE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- room_members table
CREATE TABLE room_members (
    room_id UUID REFERENCES rooms(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (room_id, user_id)
);

-- messages table
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID REFERENCES rooms(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_messages_room_created ON messages(room_id, created_at DESC);
```

### Project Structure

```
Real-time-Chat-Application/
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── config/config.go
│   │   ├── handler/
│   │   │   ├── chat.go
│   │   │   ├── room.go
│   │   │   └── user.go
│   │   ├── middleware/
│   │   │   └── auth.go
│   │   ├── model/
│   │   │   ├── message.go
│   │   │   ├── room.go
│   │   │   └── user.go
│   │   ├── repository/
│   │   │   ├── message_repo.go
│   │   │   ├── presence_repo.go
│   │   │   ├── pubsub_repo.go
│   │   │   ├── room_repo.go
│   │   │   └── user_repo.go
│   │   ├── service/
│   │   │   ├── chat_service.go
│   │   │   └── presence_service.go
│   │   └── websocket/
│   │       ├── hub.go
│   │       └── client.go
│   ├── pkg/
│   │   ├── database/postgres.go
│   │   └── redis/client.go
│   ├── migrations/
│   │   └── 001_init.sql
│   ├── go.mod
│   └── Makefile
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx
│   │   │   ├── chat/[roomId]/page.tsx
│   │   │   └── globals.css
│   │   ├── components/
│   │   │   ├── chat/
│   │   │   │   ├── ChatRoom.tsx
│   │   │   │   ├── MessageList.tsx
│   │   │   │   ├── MessageInput.tsx
│   │   │   │   └── OnlineUsers.tsx
│   │   │   └── ui/
│   │   │       ├── Button.tsx
│   │   │       └── Card.tsx
│   │   ├── hooks/
│   │   │   ├── useWebSocket.ts
│   │   │   └── useChat.ts
│   │   ├── lib/
│   │   │   └── api.ts
│   │   ├── types/index.ts
│   │   └── theme/isan-colors.ts
│   ├── package.json
│   ├── tailwind.config.ts
│   └── next.config.ts
├── docker-compose.yml
├── README.md
└── .gitignore
```

### Further Considerations

1. **Authentication** - เริ่มด้วย nickname-based ก่อน แล้วค่อยเพิ่ม JWT auth ใน Phase 2

2. **Room Management** - เริ่มด้วย multiple public rooms ก่อน แล้วค่อยเพิ่ม private rooms

3. **Message Features** - เริ่มด้วย typing indicators ผ่าน Pub/Sub ก่อน

4. **PostgreSQL Setup** - เพิ่มใน docker-compose.yml พร้อม Redis ที่มีอยู่แล้ว
