# Scalable Learning Platform (Go + React)

This repository contains a scaffolded implementation of a scalable learning platform for high-traffic events, built with a Go backend and React frontend. It follows the attached system design plan and is structured to run on AWS with Dockerized services.

## Structure

- `backend/` – Go HTTP API (chi), JWT auth scaffold, PostgreSQL database layer, student and admin endpoints, basic realtime stream.
- `frontend/` – React + Vite SPA with student and admin dashboards, using React Router and React Query.

## Database Setup

This project uses **PostgreSQL** as the primary database. The schema includes tables for:
- `users` – User accounts with roles (STUDENT, INSTRUCTOR, ADMIN)
- `sessions` – Live and recorded sessions
- `study_materials` – Course materials (PDFs, PPTs, Links)
- `exams` – Tests and assignments
- `submissions` – Student exam submissions and scores
- `notifications` – User notifications

### Running Locally with Docker

The easiest way to run the project is with Docker Compose:

```bash
docker-compose up
```

This will start:
- **PostgreSQL** on `localhost:5432` (user: postgres, password: postgres)
- **Backend API** on `http://localhost:8080`
- **Frontend** on `http://localhost:5173`

The database schema is automatically initialized on first run using `backend/schema.sql`.

### Running Locally Without Docker

#### Prerequisites
- PostgreSQL 16+
- Go 1.22+
- Node.js 18+

#### Backend Setup

1. **Create database and schema:**
```bash
createdb scalable_learning
psql scalable_learning < backend/schema.sql
```

2. **Configure environment variables:**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

3. **Run backend:**
```bash
cd backend
go mod download
go run ./cmd/api
```

The API listens on `http://localhost:8080` by default.

#### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

The app listens on `http://localhost:5173` and proxies `/api` calls to the backend.

## Database Schema

Key tables and relationships:

### Users
- Stores user information with roles (STUDENT, INSTRUCTOR, ADMIN)
- Password stored as bcrypt hash
- Unique email constraint

### Sessions
- Represents live or recorded sessions
- Status: SCHEDULED, ACTIVE, COMPLETED
- Links to created_by user

### Study Materials
- Associated with sessions
- Can be PDF, PPT, or external links
- S3 key for cloud storage integration

### Exams
- Associated with sessions
- Can be EXAM or ASSIGNMENT type
- Has open and close times

### Submissions
- Student exam submissions
- Status: SUBMITTED or GRADED
- Includes score when graded

### Notifications
- User notifications (unread flag)
- Type-based for different notification categories

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` – User registration
- `POST /api/v1/auth/login` – User login (returns JWT token)
- `GET /api/v1/auth/me` – Get current user (requires auth)

### Student Dashboard
- `GET /api/v1/student/dashboard` – Get upcoming/past sessions and notifications

### Sessions
- `GET /api/v1/sessions` – List all sessions
- `GET /api/v1/sessions/{id}` – Get session details
- `GET /api/v1/sessions/{id}/status` – Get session status
- `GET /api/v1/sessions/{id}/materials` – Get session materials

### Exams
- `POST /api/v1/exams/{id}/submissions` – Submit exam
- `GET /api/v1/exams/{id}/submissions/me` – Get student's submission

### Admin
- `POST /api/v1/admin/sessions` – Create session
- `PATCH /api/v1/admin/sessions/{id}/status` – Update session status
- `GET /api/v1/admin/sessions/{id}/stats` – Get session statistics
- `GET /api/v1/admin/submissions` – Get submission reports

## Docker Images

- `backend/Dockerfile` builds a small Go image running the API on port 8080.
- `frontend/Dockerfile` runs the development server on port 5173 (use production Dockerfile for nginx deployment).

These images are suitable for deployment on AWS ECS/EKS behind an Application Load Balancer, with CloudFront in front of the frontend image and RDS/Redis/S3 wired according to the design plan.

## AWS Infrastructure (High Level)

- **Compute**: ECS (Fargate or EC2) services for `api` and `frontend` tasks, with auto-scaling on CPU/RPS.
- **Database**: Amazon RDS (PostgreSQL) for users, sessions, materials, exams, submissions, and activity logs.
- **Cache / Realtime**: ElastiCache Redis for caching hot data and pub/sub for realtime events across API tasks.
- **Storage & CDN**:
  - S3 buckets for study materials and submission files.
  - CloudFront CDN in front of the frontend (and optionally S3 assets) for < 2s page loads globally.
- **Networking**:
  - Application Load Balancer routing `/api/*` to the Go API service and `/*` to the frontend.
  - Private subnets for ECS tasks and RDS, public subnets for ALB.
- **Operations**:
  - CloudWatch Logs and Metrics for centralized logging and alerting.

## Environment Variables

Create a `.env` file in the project root or set these variables:

```
# Backend
API_PORT=8080
JWT_SECRET=your-secret-key-change-in-production
FRONTEND_ORIGIN=http://localhost:5173
ENVIRONMENT=development

# Database (either individual params or DATABASE_URL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=scalable_learning
# OR
DATABASE_URL=postgres://user:password@host:5432/dbname?sslmode=disable
```
  - Rolling deployments via ECS for zero-downtime releases.


