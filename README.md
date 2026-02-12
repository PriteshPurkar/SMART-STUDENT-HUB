# Scalable Learning Platform (Go + React)

This repository contains a scaffolded implementation of a scalable learning platform for high-traffic events, built with a Go backend and React frontend. It follows the attached system design plan and is structured to run on AWS with Dockerized services.

## Structure

- `backend/` – Go HTTP API (chi), JWT auth scaffold, student and admin endpoints, basic realtime stream.
- `frontend/` – React + Vite SPA with student and admin dashboards, using React Router and React Query.

## Running Locally

### Backend

```bash
cd backend
go run ./cmd/api
```

The API listens on `http://localhost:8080` by default.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

The app listens on `http://localhost:5173` and proxies `/api` calls to the backend.

## Docker Images

- `backend/Dockerfile` builds a small Go image running the API on port 8080.
- `frontend/Dockerfile` builds static assets and serves them via Nginx on port 80.

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
  - Rolling deployments via ECS for zero-downtime releases.


