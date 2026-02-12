# Architecture Diagram

## System Architecture After PostgreSQL Integration

```
┌─────────────────────────────────────────────────────────────────────┐
│                          CLIENT LAYER                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│   ┌──────────────────────────────────────────────────────────────┐  │
│   │                   Frontend (React + Vite)                   │  │
│   │              http://localhost:5173                          │  │
│   │                                                              │  │
│   │   ├── Student Dashboard                                    │  │
│   │   ├── Faculty Dashboard                                    │  │
│   │   ├── Admin Dashboard                                      │  │
│   │   └── Authentication Pages                                 │  │
│   └───────────────────────────┬────────────────────────────────┘  │
│                               │                                     │
│                               ↓ HTTP REST API                       │
└───────────────────────────────┼─────────────────────────────────────┘

┌───────────────────────────────┬─────────────────────────────────────┐
│                    API LAYER (Backend)                              │
├───────────────────────────────┼─────────────────────────────────────┤
│                               ↓                                      │
│   ┌──────────────────────────────────────────────────────────────┐  │
│   │                  Go HTTP Server (Chi)                        │  │
│   │              http://localhost:8080                          │  │
│   │                                                              │  │
│   │  ┌────────────────────────────────────────────────────────┐ │  │
│   │  │            Middleware Layer                           │ │  │
│   │  ├─ CORS Handler                                         │ │  │
│   │  ├─ JWT Authentication                                   │ │  │
│   │  └─ Request Logging                                      │ │  │
│   │  └────────────────────────────────────────────────────────┘ │  │
│   │                          │                                   │  │
│   │                          ↓                                   │  │
│   │  ┌────────────────────────────────────────────────────────┐ │  │
│   │  │            HTTP API Handlers                          │ │  │
│   │  ├─ Auth Handler (register, login, me)                   │ │  │
│   │  ├─ Student Handler (dashboard)                          │ │  │
│   │  ├─ Session Handler (list, get, materials)               │ │  │
│   │  ├─ Exam Handler (submit, get submissions)               │ │  │
│   │  └─ Admin Handler (manage sessions, reports)             │ │  │
│   │  └────────────────────────────────────────────────────────┘ │  │
│   │                          │                                   │  │
│   │                          ↓                                   │  │
│   │  ┌────────────────────────────────────────────────────────┐ │  │
│   │  │            Service Layer                              │ │  │
│   │  ├─ UserService (auth, lookup)                           │ │  │
│   │  ├─ SessionService (CRUD, queries)                       │ │  │
│   │  ├─ ExamService (exams, submissions)                     │ │  │
│   │  └─ NotificationService (notifications)                  │ │  │
│   │  └────────────────────────────────────────────────────────┘ │  │
│   │                          │                                   │  │
│   │                          ↓                                   │  │
│   │  ┌────────────────────────────────────────────────────────┐ │  │
│   │  │         Database Connection Pool                      │ │  │
│   │  ├─ lib/pq Driver                                        │ │  │
│   │  ├─ Connection Pooling (25 max, 5 idle)                  │ │  │
│   │  ├─ Prepared Statements                                  │ │  │
│   │  └─ Error Handling                                       │ │  │
│   │  └────────────────────────────────────────────────────────┘ │  │
│   └────────────────────────┬─────────────────────────────────────┘  │
└────────────────────────────┼──────────────────────────────────────────┘
                             │
                             ↓ TCP/IP Port 5432
                             
┌─────────────────────────────────────────────────────────────────────┐
│                    DATABASE LAYER                                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│   ┌──────────────────────────────────────────────────────────────┐  │
│   │                PostgreSQL 16                                │  │
│   │              localhost:5432                                 │  │
│   │                                                              │  │
│   │  Tables:                                                    │  │
│   │  ┌──────────────────────────────────────────────────────┐  │  │
│   │  │ users                                                │  │  │
│   │  │ ├─ id (BIGSERIAL PK)                                │  │  │
│   │  │ ├─ name, email, password_hash                       │  │  │
│   │  │ ├─ role (ENUM: STUDENT, INSTRUCTOR, ADMIN)          │  │  │
│   │  │ ├─ created_at, updated_at                           │  │  │
│   │  │ └─ indices: email (UNIQUE), created_at              │  │  │
│   │  └──────────────────────────────────────────────────────┘  │  │
│   │                                                              │  │
│   │  ┌──────────────────────────────────────────────────────┐  │  │
│   │  │ sessions                                             │  │  │
│   │  │ ├─ id (BIGSERIAL PK)                                │  │  │
│   │  │ ├─ title, description, start_time, end_time         │  │  │
│   │  │ ├─ status (ENUM: SCHEDULED, ACTIVE, COMPLETED)      │  │  │
│   │  │ ├─ video_url, created_by (FK → users)              │  │  │
│   │  │ └─ indices: status, start_time                       │  │  │
│   │  └──────────────────────────────────────────────────────┘  │  │
│   │                                                              │  │
│   │  ┌──────────────────────────────────────────────────────┐  │  │
│   │  │ study_materials                                      │  │  │
│   │  │ ├─ id (BIGSERIAL PK)                                │  │  │
│   │  │ ├─ session_id (FK → sessions)                       │  │  │
│   │  │ ├─ title, type (ENUM: PDF, PPT, LINK)              │  │  │
│   │  │ ├─ s3_key, url, uploaded_by (FK → users)           │  │  │
│   │  │ └─ indices: session_id                              │  │  │
│   │  └──────────────────────────────────────────────────────┘  │  │
│   │                                                              │  │
│   │  ┌──────────────────────────────────────────────────────┐  │  │
│   │  │ exams                                                │  │  │
│   │  │ ├─ id (BIGSERIAL PK)                                │  │  │
│   │  │ ├─ session_id (FK → sessions)                       │  │  │
│   │  │ ├─ title, open_time, close_time                     │  │  │
│   │  │ ├─ max_score, type (ENUM: EXAM, ASSIGNMENT)         │  │  │
│   │  │ └─ indices: session_id                              │  │  │
│   │  └──────────────────────────────────────────────────────┘  │  │
│   │                                                              │  │
│   │  ┌──────────────────────────────────────────────────────┐  │  │
│   │  │ submissions                                          │  │  │
│   │  │ ├─ id (BIGSERIAL PK)                                │  │  │
│   │  │ ├─ exam_id (FK → exams)                             │  │  │
│   │  │ ├─ student_id (FK → users)                          │  │  │
│   │  │ ├─ submitted_at, file_s3_key                        │  │  │
│   │  │ ├─ score, status (ENUM: SUBMITTED, GRADED)          │  │  │
│   │  │ └─ indices: exam_id + student_id, status            │  │  │
│   │  └──────────────────────────────────────────────────────┘  │  │
│   │                                                              │  │
│   │  ┌──────────────────────────────────────────────────────┐  │  │
│   │  │ notifications                                        │  │  │
│   │  │ ├─ id (BIGSERIAL PK)                                │  │  │
│   │  │ ├─ user_id (FK → users)                             │  │  │
│   │  │ ├─ type, message, is_read                           │  │  │
│   │  │ ├─ created_at                                       │  │  │
│   │  │ └─ indices: user_id, is_read                        │  │  │
│   │  └──────────────────────────────────────────────────────┘  │  │
│   │                                                              │  │
│   └──────────────────────────────────────────────────────────────┘  │
│                                                                       │
│   Data Persistence:                                                 │
│   ├─ PostgreSQL Volume (persistent across restarts)                 │
│   └─ Automatic backups (via pg_dump)                                │
│                                                                       │
└─────────────────────────────────────────────────────────────────────┘
```

## Docker Compose Services

```
┌────────────────────────────────────────────────────────────────┐
│                    docker-compose.yml                          │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Network: scalable-learning-network (bridge)                  │
│                                                                 │
│  Services:                                                     │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │ postgres                                                │  │
│  │ ├─ Image: postgres:16-alpine                           │  │
│  │ ├─ Port: 5432                                          │  │
│  │ ├─ Volume: postgres_data (persistent)                  │  │
│  │ ├─ Init Script: schema.sql                             │  │
│  │ └─ Health Check: pg_isready                            │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │ backend                                                 │  │
│  │ ├─ Build: ./backend/Dockerfile                         │  │
│  │ ├─ Port: 8080                                          │  │
│  │ ├─ Depends On: postgres (health check)                 │  │
│  │ ├─ Env: DB_HOST=postgres, DB_PORT=5432                │  │
│  │ └─ Network: scalable-learning-network                  │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │ frontend                                                │  │
│  │ ├─ Build: ./frontend/Dockerfile                        │  │
│  │ ├─ Port: 5173                                          │  │
│  │ ├─ Depends On: backend                                 │  │
│  │ ├─ Env: VITE_API_URL=http://localhost:8080/api/v1     │  │
│  │ └─ Network: scalable-learning-network                  │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

## Data Flow

### User Registration Flow
```
1. Frontend Form Input
   ↓
2. POST /api/v1/auth/register (JSON: name, email, password)
   ↓
3. Router → Auth Handler
   ↓
4. Validate Input
   ↓
5. UserService.CreateUser()
   ↓
6. Hash Password (bcrypt)
   ↓
7. INSERT INTO users
   ↓
8. Return User Object (200 Created)
   ↓
9. Frontend: User Created ✓
```

### Login Flow with Database
```
1. Frontend Form Input
   ↓
2. POST /api/v1/auth/login (JSON: email, password)
   ↓
3. Router → Auth Handler
   ↓
4. UserService.GetUserByEmail(email)
   ↓
5. SELECT * FROM users WHERE email = ?
   ↓
6. UserService.VerifyPassword()
   ↓
7. Compare bcrypt hashes
   ↓
8. Generate JWT Token
   ↓
9. Return Token + User (200 OK)
   ↓
10. Frontend: Login Successful ✓
```

### Session Query Flow (Protected)
```
1. GET /api/v1/sessions (with JWT in Authorization header)
   ↓
2. Router → JWT Middleware
   ↓
3. Validate Token
   ↓
4. Extract UserID from Token
   ↓
5. Router → Session Handler
   ↓
6. SessionService.GetAllSessions()
   ↓
7. SELECT * FROM sessions ORDER BY start_time DESC
   ↓
8. Build Response
   ↓
9. Return Sessions (200 OK)
   ↓
10. Frontend: Display Sessions ✓
```

## Security Flow

```
Password Management:
Input Password → bcrypt.GenerateFromPassword() → Hash Stored
Login Password → bcrypt.CompareHashAndPassword() → Compare Hash

SQL Injection Prevention:
User Input → Prepared Statement (parameterized query)
SELECT * FROM users WHERE email = $1 (postgres/pq)

Authentication:
JWT Token → Extract UserID & Role → Verify in Middleware
Authorization: Role-based access control in handlers
```

## Performance Optimizations

```
Database Level:
├─ Connection Pooling (25 max, 5 idle)
├─ Prepared Statements
├─ Indexed Columns
│  ├─ users: email (UNIQUE), created_at
│  ├─ sessions: status, start_time
│  ├─ exams: session_id
│  ├─ submissions: exam_id + student_id, status
│  └─ notifications: user_id, is_read
└─ Foreign Key Constraints

Application Level:
├─ Connection Reuse
├─ Lazy Loading
├─ Query Optimization
└─ Error Handling
```
