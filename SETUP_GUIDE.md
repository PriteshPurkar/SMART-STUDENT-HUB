# PostgreSQL Setup Guide

This guide provides step-by-step instructions to set up PostgreSQL for the Scalable Learning Platform.

## Quick Start with Docker Compose

The fastest way to get started:

```bash
docker-compose up
```

This command will:
1. Start PostgreSQL on `localhost:5432`
2. Initialize the database schema from `backend/schema.sql`
3. Start the backend API on `localhost:8080`
4. Start the frontend on `localhost:5173`

All services are automatically connected and ready to use.

## Local Setup Without Docker

### Prerequisites
- PostgreSQL 16+ ([Download](https://www.postgresql.org/download/))
- Go 1.22+ ([Download](https://golang.org/dl/))
- Node.js 18+ ([Download](https://nodejs.org/))

### Step 1: Install PostgreSQL

**macOS (using Homebrew):**
```bash
brew install postgresql
brew services start postgresql
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql
```

**Windows:**
Download installer from [postgresql.org](https://www.postgresql.org/download/windows/) and follow the installation wizard.

### Step 2: Create Database and User

```bash
# Connect to PostgreSQL (default user: postgres)
psql postgres

# Create database
CREATE DATABASE scalable_learning;

# Create user (or use postgres user)
CREATE USER scalable_user WITH PASSWORD 'your_secure_password';

# Grant privileges
GRANT ALL PRIVILEGES ON DATABASE scalable_learning TO scalable_user;

# Connect to the database
\c scalable_learning

# Exit psql
\q
```

### Step 3: Initialize Schema

```bash
# Navigate to project directory
cd /path/to/project

# Load schema
psql -U scalable_user -d scalable_learning -f backend/schema.sql
```

You should see output like:
```
CREATE TYPE
CREATE TYPE
CREATE TYPE
CREATE TYPE
CREATE TYPE
CREATE TABLE
CREATE INDEX
...
```

### Step 4: Configure Environment

```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your settings
```

Example `.env` file:
```
# Backend
API_PORT=8080
JWT_SECRET=dev-secret-key-change-in-production
FRONTEND_ORIGIN=http://localhost:5173
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=scalable_user
DB_PASSWORD=your_secure_password
DB_NAME=scalable_learning
```

### Step 5: Run Backend

```bash
cd backend

# Download dependencies
go mod download

# Run the server
go run ./cmd/api
```

You should see:
```
API listening on :8080
Environment: development
Database: postgres://***:***@***:***/**?
```

### Step 6: Run Frontend

In a new terminal:

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend should be available at `http://localhost:5173`

## Database Schema Overview

### Tables Created

| Table | Purpose |
|-------|---------|
| `users` | User accounts with roles |
| `sessions` | Live and recorded sessions |
| `study_materials` | Course materials |
| `exams` | Tests and assignments |
| `submissions` | Student submissions |
| `notifications` | User notifications |

### Relationships

```
users (1) ──→ (many) sessions (created_by)
users (1) ──→ (many) submissions (student_id)
sessions (1) ──→ (many) study_materials
sessions (1) ──→ (many) exams
exams (1) ──→ (many) submissions
```

## Verifying Setup

### Test Backend Connection

```bash
# Health check
curl http://localhost:8080/api/v1/healthz

# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Test Database Connection

```bash
# Connect directly to database
psql -U scalable_user -d scalable_learning

# Check tables
\dt

# Check rows in users table
SELECT * FROM users;

# Exit
\q
```

## Troubleshooting

### "connection refused" error

**Issue:** Backend can't connect to PostgreSQL

**Solutions:**
1. Verify PostgreSQL is running: `pg_isready -h localhost -p 5432`
2. Check database credentials in `.env`
3. Ensure database exists: `psql -U postgres -l | grep scalable_learning`
4. Check firewall settings on port 5432

### "database does not exist" error

```bash
# Create the database
createdb -U postgres scalable_learning

# Initialize schema
psql -U postgres -d scalable_learning -f backend/schema.sql
```

### "permission denied" error

```bash
# Grant privileges to user
psql -U postgres -d scalable_learning -c \
  "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO scalable_user;"
```

### Migrations Failed in Docker

If schema initialization fails in Docker, you can manually apply it:

```bash
# Connect to the running postgres container
docker exec -it scalable-learning-postgres psql -U postgres -d scalable_learning

# From inside psql:
\i /docker-entrypoint-initdb.d/schema.sql

# Exit
\q
```

## Backup and Restore

### Backup Database

```bash
pg_dump -U scalable_user scalable_learning > backup.sql
```

### Restore Database

```bash
psql -U scalable_user scalable_learning < backup.sql
```

## Production Considerations

1. **Security**: 
   - Use strong passwords for database user
   - Never commit `.env` files with real credentials
   - Use environment variables in production
   - Enable SSL/TLS for database connections

2. **Performance**:
   - Add indices for frequently queried columns (done in schema)
   - Configure connection pooling (pgBouncer)
   - Regular maintenance (VACUUM, ANALYZE)

3. **Backups**:
   - Set up automated daily backups
   - Test restore procedures regularly
   - Use AWS RDS for managed backups

4. **Monitoring**:
   - Monitor database size and connections
   - Set up alerts for unusual activity
   - Log slow queries
   - Monitor replication lag

## Additional Resources

- [PostgreSQL Official Documentation](https://www.postgresql.org/docs/)
- [GoSQL Driver for PostgreSQL](https://github.com/lib/pq)
- [Docker PostgreSQL Image](https://hub.docker.com/_/postgres)
- [AWS RDS PostgreSQL](https://aws.amazon.com/rds/postgresql/)
