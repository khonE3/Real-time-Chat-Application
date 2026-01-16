# 🏯 แชทอีสาน - หนองบัวลำภู | Isan Real-time Chat Application

แอปพลิเคชันแชทเรียลไทม์ ธีมอีสานหนองบัวลำภู สร้างด้วย Go Fiber, Redis 7.4, PostgreSQL และ Next.js 16

![Isan Chat Banner](https://img.shields.io/badge/ธีม-อีสานหนองบัวลำภู-gold?style=for-the-badge)
![Go Version](https://img.shields.io/badge/Go-1.23-blue?style=flat-square)
![Next.js](https://img.shields.io/badge/Next.js-16.1-black?style=flat-square)
![Redis](https://img.shields.io/badge/Redis-7.4-red?style=flat-square)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17.2-blue?style=flat-square)

## ✨ Features

### Core Features (MVP)
- ✅ **Real-time Messaging** - ส่งข้อความเรียลไทม์ผ่าน WebSocket + Redis Pub/Sub
- ✅ **Message History** - เก็บประวัติข้อความใน PostgreSQL + Redis Streams
- ✅ **Online User Tracking** - ติดตามผู้ใช้ออนไลน์ด้วย Redis Sorted Sets
- ✅ **Multiple Chat Rooms** - หลายห้องแชท (public)
- ✅ **Isan Theme UI** - ธีมสีทองวัด, แดงผ้าไหม, เขียวทุ่งนา
- ✅ **Simple Auth** - ระบบ nickname-based authentication
- ✅ **Typing Indicators** - แสดงสถานะกำลังพิมพ์
- ✅ **Unread Count** - นับข้อความที่ยังไม่ได้อ่าน

## 🛠️ Tech Stack

| Technology | Version | Purpose |
|------------|---------|---------|  
| **Go** | 1.23 | Backend runtime |
| **Go Fiber** | v2.52.5 | Web framework + WebSocket |
| **go-redis** | v9.7.2 | Redis client |
| **pgx** | v5.7.2 | PostgreSQL driver |
| **Next.js** | 16.1 | Frontend framework |
| **React** | 19.2 | UI library |
| **TailwindCSS** | 4.1 | Styling |
| **Bun** | 1.3.3 | Package manager |
| **Redis** | 7.4 | Real-time (Pub/Sub, Streams, Presence) |
| **PostgreSQL** | 17.2 | Persistent storage |

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLIENT (Next.js)                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │  WebSocket  │  │  REST API   │  │  React UI   │              │
│  └──────┬──────┘  └──────┬──────┘  └─────────────┘              │
└─────────┼────────────────┼──────────────────────────────────────┘
          │                │
          ▼                ▼
┌─────────────────────────────────────────────────────────────────┐
│                       SERVER (Go Fiber)                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │  WebSocket  │  │  Handlers   │  │  Services   │              │
│  │    Hub      │  │             │  │             │              │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘              │
└─────────┼────────────────┼────────────────┼─────────────────────┘
          │                │                │
          ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────────┐
│                        DATA LAYER                                │
│  ┌──────────────────────┐    ┌──────────────────────┐           │
│  │      PostgreSQL      │    │        Redis         │           │
│  │   (Persistent Data)  │    │   (Real-time Data)   │           │
│  ├──────────────────────┤    ├──────────────────────┤           │
│  │ • Users              │    │ • Pub/Sub (live msg) │           │
│  │ • Rooms              │    │ • Streams (recent)   │           │
│  │ • Messages (history) │    │ • Sorted Sets        │           │
│  │ • Room members       │    │   (online users)     │           │
│  └──────────────────────┘    └──────────────────────┘           │
└─────────────────────────────────────────────────────────────────┘
```

## 📁 Project Structure

```
Real-time-Chat-Application/
├── backend/                    # Go Fiber Backend
│   ├── cmd/server/main.go     # Entry point
│   ├── internal/
│   │   ├── config/            # Configuration
│   │   ├── handler/           # HTTP handlers
│   │   ├── middleware/        # Middleware (auth, cors)
│   │   ├── model/             # Data models
│   │   ├── repository/        # Database operations
│   │   ├── service/           # Business logic
│   │   └── websocket/         # WebSocket hub
│   ├── pkg/
│   │   ├── database/          # PostgreSQL client
│   │   └── redis/             # Redis client
│   ├── migrations/            # SQL migrations
│   ├── go.mod
│   └── Makefile
├── frontend/                   # Next.js Frontend
│   ├── src/
│   │   ├── app/               # Pages (App Router)
│   │   ├── components/        # React components
│   │   ├── hooks/             # Custom hooks
│   │   ├── lib/               # Utilities
│   │   ├── types/             # TypeScript types
│   │   └── theme/             # Isan theme colors
│   ├── package.json
│   └── next.config.ts
├── docker-compose.yml          # Docker services
├── README.md
└── .gitignore
```

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- Bun 1.3+
- Docker & Docker Compose (for Redis)
- PostgreSQL 17 (via DBeaver or Docker)
- Redis 7.4 (via Docker)

### 1. Clone the repository

```bash
git clone https://github.com/khonE3/Real-time-Chat-Application.git
cd Real-time-Chat-Application
```

### 2. Start Redis with Docker

```bash
docker-compose up -d redis
# or run Redis only:
docker run -d --name chat_redis -p 6379:6379 redis:7.4-alpine
```

### 3. Setup PostgreSQL Database

Using DBeaver or pgAdmin:
1. Connect to PostgreSQL server
2. Create database: `CREATE DATABASE chatdb;`
3. Run migration script from `backend/migrations/001_init.sql`

### 4. Setup Backend

```bash
cd backend

# Copy environment file
cp .env.example .env

# Edit .env with your PostgreSQL password
# DB_PASSWORD=your_password

# Download dependencies
go mod tidy

# Run the server
go run cmd/server/main.go
```

Backend will start at `http://localhost:3001`

### 5. Setup Frontend

```bash
cd frontend

# Install dependencies with Bun
bun install

# Run development server
bun dev
```

Frontend will start at `http://localhost:3000`

### 6. Open the application

Visit `http://localhost:3000` in your browser 🎉

## 🔌 API Endpoints

### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/users` | Create/Get user |
| GET | `/api/users/:id` | Get user by ID |
| GET | `/api/users/username/:username` | Get user by username |

### Rooms

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/rooms` | List all rooms (add ?userId=xxx for unread counts) |
| POST | `/api/rooms` | Create new room |
| GET | `/api/rooms/:id` | Get room by ID |
| POST | `/api/rooms/:id/join` | Join a room |
| GET | `/api/rooms/:id/members` | Get room members |
| POST | `/api/rooms/:id/read?userId=xxx` | Mark messages as read |
| GET | `/api/rooms/:id/unread?userId=xxx` | Get unread message count |

### Messages

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/rooms/:id/messages` | Get room messages |

### WebSocket

| Endpoint | Description |
|----------|-------------|
| `ws://localhost:3001/ws/:roomId?userId=xxx&username=xxx&displayName=xxx` | WebSocket connection |

#### WebSocket Message Types

```typescript
// Outgoing (Client → Server)
{ type: "message", content: "Hello!" }
{ type: "typing" }
{ type: "stop_typing" }

// Incoming (Server → Client)
{ type: "message", payload: Message }
{ type: "history", payload: Message[] }
{ type: "online_users", payload: OnlineUser[] }
{ type: "typing", payload: TypingUser }
{ type: "stop_typing", payload: TypingUser }
{ type: "presence", payload: PresencePayload }
```

## 🎨 Isan Theme Colors

| Color | Hex | Inspiration |
|-------|-----|-------------|
| 🟡 **Gold** | `#D4A12A` | สีทองวัด - Temple Gold |
| 🔴 **Silk** | `#DC143C` | สีแดงผ้าไหม - Traditional Silk |
| 🟢 **Paddy** | `#228B22` | สีเขียวทุ่งนา - Rice Paddy Green |
| 🟤 **Earth** | `#8B4513` | สีน้ำตาลดิน - Laterite Soil |
| 🟡 **Cotton** | `#FFF8DC` | สีผ้าฝ้าย - Handwoven Cotton |

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Rooms Table
```sql
CREATE TABLE rooms (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_private BOOLEAN DEFAULT FALSE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP
);
```

### Messages Table
```sql
CREATE TABLE messages (
    id UUID PRIMARY KEY,
    room_id UUID REFERENCES rooms(id),
    user_id UUID REFERENCES users(id),
    content TEXT NOT NULL,
    message_type VARCHAR(20),
    created_at TIMESTAMP
);
```

## 🔧 Environment Variables

### Backend (.env)
```bash
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=3001

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=chatdb
DB_SSLMODE=disable

# Redis
REDIS_URL=localhost:6379
REDIS_PASSWORD=

# CORS
CORS_ORIGINS=http://localhost:3000
```

### Frontend (.env.local)
```bash
NEXT_PUBLIC_API_URL=http://localhost:3001
NEXT_PUBLIC_WS_URL=ws://localhost:3001
```

## 📝 License

MIT License - สร้างด้วย ❤️ จากดินแดนอีสาน

## 🙏 Acknowledgments

- Go Fiber Team
- Redis Community
- Next.js Team
- TailwindCSS Team
- หนองบัวลำภู - บ้านเกิดของความคิดสร้างสรรค์ 🎋
