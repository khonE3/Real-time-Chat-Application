# 🎋 แชทอีสาน - Real-time Chat Application

แอปพลิเคชันแชทแบบเรียลไทม์ พัฒนาด้วย Go Fiber, PostgreSQL, Redis และ Next.js

![Thai](https://img.shields.io/badge/Language-Thai-blue)
![Go](https://img.shields.io/badge/Backend-Go%20Fiber-00ADD8)
![Next.js](https://img.shields.io/badge/Frontend-Next.js%2016-black)
![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-336791)
![Redis](https://img.shields.io/badge/Cache-Redis-DC382D)

## ✨ Features

- 💬 **Real-time Messaging** - ส่งข้อความทันทีผ่าน WebSocket
- 👥 **Online Users** - เห็นว่าใครออนไลน์อยู่ในห้อง
- ✍️ **Typing Indicator** - รู้ว่าใครกำลังพิมพ์ข้อความ
- 🔔 **Unread Count** - แจ้งเตือนจำนวนข้อความที่ยังไม่ได้อ่าน
- 🏠 **Multiple Rooms** - หลายห้องแชท แยกหัวข้อสนทนา
- 📱 **Responsive Design** - ใช้งานได้ทุกอุปกรณ์
- 🎨 **Isan Theme** - ธีมสีสันแบบอีสาน สวยงามเป็นเอกลักษณ์

## 🛠️ Tech Stack

### Backend
- **Go 1.21+** - ภาษาโปรแกรมหลัก
- **Fiber v2** - Web framework ประสิทธิภาพสูง
- **PostgreSQL** - ฐานข้อมูลหลัก
- **Redis** - Cache, Pub/Sub, Presence tracking
- **WebSocket** - การสื่อสารแบบ real-time

### Frontend
- **Next.js 16** - React framework
- **TypeScript** - Type-safe JavaScript
- **Tailwind CSS 4** - Utility-first CSS
- **Bun** - JavaScript runtime & package manager

## 📁 Project Structure

```
Real-time-Chat-Application/
├── backend/
│   ├── cmd/server/          # Entry point
│   ├── internal/
│   │   ├── config/          # Configuration
│   │   ├── handler/         # HTTP handlers
│   │   ├── middleware/      # Middlewares
│   │   ├── model/           # Data models
│   │   ├── repository/      # Database layer
│   │   ├── service/         # Business logic
│   │   └── websocket/       # WebSocket hub
│   ├── migrations/          # SQL migrations
│   ├── pkg/
│   │   ├── database/        # PostgreSQL client
│   │   └── redis/           # Redis client
│   ├── .env                 # Environment variables
│   └── go.mod
│
├── frontend/
│   ├── src/
│   │   ├── app/             # Next.js pages
│   │   ├── components/      # React components
│   │   ├── hooks/           # Custom hooks
│   │   ├── lib/             # Utilities
│   │   └── types/           # TypeScript types
│   ├── package.json
│   └── next.config.ts
│
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+ หรือ Bun 1.0+
- PostgreSQL 14+
- Redis 7+

### 1. Clone Repository

```bash
git clone https://github.com/khonE3/Real-time-Chat-Application.git
cd Real-time-Chat-Application
```

### 2. Setup Database

สร้างฐานข้อมูล PostgreSQL:

```sql
CREATE DATABASE chatdb;
```

รัน migration:

```bash
psql -U postgres -d chatdb -f backend/migrations/001_init.sql
```

### 3. Setup Backend

```bash
cd backend

# Copy environment file
cp .env.example .env

# Edit .env with your settings
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=your_password
# DB_NAME=chatdb
# REDIS_URL=localhost:6379

# Run backend
go run ./cmd/server
```

Backend จะรันที่ `http://localhost:3001`

### 4. Setup Frontend

```bash
cd frontend

# Install dependencies (ใช้ bun หรือ npm)
bun install
# หรือ
npm install

# Run development server
bun dev
# หรือ
npm run dev
```

Frontend จะรันที่ `http://localhost:3000`

### 5. Start Redis

```bash
# Docker
docker run -d --name chat_redis -p 6379:6379 redis:7.4-alpine

# หรือ Redis ที่ติดตั้งในเครื่อง
redis-server
```

## 🔧 Environment Variables

### Backend (.env)

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_HOST` | Server host | `0.0.0.0` |
| `SERVER_PORT` | Server port | `3001` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL user | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | - |
| `DB_NAME` | Database name | `chatdb` |
| `REDIS_URL` | Redis URL | `localhost:6379` |
| `CORS_ORIGINS` | Allowed origins | `http://localhost:3000` |

### Frontend

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | `http://127.0.0.1:3001` |
| `NEXT_PUBLIC_WS_URL` | WebSocket URL | `ws://127.0.0.1:3001` |

## 📡 API Endpoints

### Users
- `POST /api/users` - สร้าง/ล็อกอิน user
- `GET /api/users/:id` - ดึงข้อมูล user
- `GET /api/users/username/:username` - ค้นหา user จาก username

### Rooms
- `GET /api/rooms` - รายการห้องแชททั้งหมด
- `POST /api/rooms` - สร้างห้องใหม่
- `GET /api/rooms/:id` - ดึงข้อมูลห้อง
- `POST /api/rooms/:id/join` - เข้าร่วมห้อง
- `GET /api/rooms/:id/members` - รายการสมาชิกในห้อง
- `POST /api/rooms/:id/read` - อ่านข้อความแล้ว
- `GET /api/rooms/:id/unread` - จำนวนข้อความที่ยังไม่อ่าน
- `GET /api/rooms/:id/messages` - ประวัติข้อความ

### WebSocket
- `WS /ws/:roomId?userId=...&username=...&displayName=...`

#### WebSocket Message Types

**Incoming (Client → Server):**
```json
{ "type": "message", "content": "Hello!" }
{ "type": "typing" }
{ "type": "stop_typing" }
```

**Outgoing (Server → Client):**
```json
{ "type": "message", "payload": { ... } }
{ "type": "history", "payload": [ ... ] }
{ "type": "online_users", "payload": [ ... ] }
{ "type": "typing", "payload": { ... } }
{ "type": "presence", "payload": { ... } }
{ "type": "error", "payload": "Error message" }
```

## 🎨 Screenshots

### หน้าหลัก
- แสดงรายการห้องแชท
- จำนวนคนออนไลน์
- Unread count badge

### ห้องแชท
- ข้อความแบบ real-time
- Typing indicator
- รายการคนออนไลน์
- สถานะการเชื่อมต่อ

## 🐛 Troubleshooting

### WebSocket ไม่เชื่อมต่อ (Windows)

ถ้าใช้ Windows และ WebSocket ไม่เชื่อมต่อ ให้ตรวจสอบว่า:
1. ใช้ `127.0.0.1` แทน `localhost` (ปัญหา IPv6)
2. Redis กำลังทำงานอยู่
3. Backend รันอยู่ที่ port 3001

### Database Connection Error

ตรวจสอบ:
1. PostgreSQL กำลังทำงาน
2. ข้อมูลการเชื่อมต่อใน `.env` ถูกต้อง
3. ถ้ารหัสผ่านมีอักขระพิเศษ (เช่น `%`) ต้อง escape ให้ถูกต้อง

## 📄 License

MIT License - ใช้งานได้อิสระ

## 👨‍💻 Author

**khonE3**

---

🏯 *"พูดคุยแลกเปลี่ยน เหมือนนั่งกินข้าวเหนียวริมโขง"* 🎋
