# Quick Reference Guide

## 📚 Documentation Files

| File | Purpose |
|------|---------|
| [README.md](README.md) | Main documentation and API reference |
| [SETUP_GUIDE.md](SETUP_GUIDE.md) | Step-by-step local setup instructions |
| [CHANGES_SUMMARY.md](CHANGES_SUMMARY.md) | Summary of all changes made |
| [POSTGRESQL_INTEGRATION.md](POSTGRESQL_INTEGRATION.md) | Technical PostgreSQL integration details |
| [ARCHITECTURE.md](ARCHITECTURE.md) | System architecture and data flow diagrams |
| [.env.example](.env.example) | Environment variables template |

## 🚀 Quick Start

### Docker (Fastest)
```bash
docker-compose up
```
Then access:
- Frontend: http://localhost:5173
- Backend: http://localhost:8080
- Database: localhost:5432 (postgres/postgres)

### Local Development
```bash
# 1. Setup database
createdb scalable_learning
psql scalable_learning < backend/schema.sql

# 2. Setup environment
cp .env.example .env

# 3. Terminal 1: Backend
cd backend && go run ./cmd/api

# 4. Terminal 2: Frontend
cd frontend && npm install && npm run dev
```

## 🔧 Utility Scripts

| Script | Purpose | Usage |
|--------|---------|-------|
| `db.sh` | Database management | `./db.sh [init\|reset\|backup\|restore]` |
| `start.sh` | Quick start (Linux/Mac) | `./start.sh` |
| `start.bat` | Quick start (Windows) | `start.bat` |

## 📝 API Endpoints Reference

### Auth
```
POST /api/v1/auth/register
POST /api/v1/auth/login
GET /api/v1/auth/me
```

### Student
```
GET /api/v1/student/dashboard
```

### Sessions
```
GET /api/v1/sessions
GET /api/v1/sessions/{id}
GET /api/v1/sessions/{id}/status
GET /api/v1/sessions/{id}/materials
```

### Exams
```
POST /api/v1/exams/{id}/submissions
GET /api/v1/exams/{id}/submissions/me
```

### Admin
```
POST /api/v1/admin/sessions
PATCH /api/v1/admin/sessions/{id}/status
GET /api/v1/admin/sessions/{id}/stats
GET /api/v1/admin/submissions
```

## 💾 Database Schema

### Tables
- `users` - User accounts with roles
- `sessions` - Live/recorded sessions
- `study_materials` - Course materials
- `exams` - Tests and assignments
- `submissions` - Student work and scores
- `notifications` - User notifications

### Key Columns
```sql
users.id (PK)
users.email (UNIQUE)
users.role (ENUM: STUDENT, INSTRUCTOR, ADMIN)

sessions.id (PK)
sessions.status (ENUM: SCHEDULED, ACTIVE, COMPLETED)
sessions.created_by (FK: users.id)

exams.id (PK)
exams.session_id (FK: sessions.id)

submissions.id (PK)
submissions.exam_id (FK: exams.id)
submissions.student_id (FK: users.id)
```

## 🌍 Environment Variables

```
# Backend
API_PORT=8080
JWT_SECRET=your-secret-key
FRONTEND_ORIGIN=http://localhost:5173
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=scalable_learning
```

## 🐛 Troubleshooting

### Database Connection Failed
```bash
# Check if PostgreSQL is running
pg_isready -h localhost -p 5432

# Check database exists
psql -U postgres -l | grep scalable_learning

# Reinitialize
createdb scalable_learning
psql scalable_learning < backend/schema.sql
```

### Port Already in Use
```bash
# Change port in .env
API_PORT=8081

# Or kill process using port
lsof -ti:8080 | xargs kill -9  # Linux/Mac
netstat -ano | findstr :8080   # Windows
```

### Docker Issues
```bash
# Clean up containers
docker-compose down

# Remove volumes (careful!)
docker-compose down -v

# Rebuild
docker-compose up --build
```

## 📂 Project Structure

```
Project/
├── backend/
│   ├── cmd/
│   │   └── api/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── db/              ← Database services
│   │   ├── httpapi/         ← API handlers
│   │   ├── middleware/
│   │   ├── models/
│   │   └── realtime/
│   ├── schema.sql           ← Database schema
│   ├── go.mod & go.sum
│   └── Dockerfile
├── frontend/
│   ├── src/
│   ├── package.json
│   └── Dockerfile
├── docker-compose.yml
├── .env.example
├── README.md
├── SETUP_GUIDE.md
├── CHANGES_SUMMARY.md
├── POSTGRESQL_INTEGRATION.md
├── ARCHITECTURE.md
├── db.sh
├── start.sh & start.bat
└── This file (QUICK_REFERENCE.md)
```

## ✅ Verification Checklist

After setup, verify:
- [ ] `docker-compose up` completes without errors
- [ ] Frontend loads at http://localhost:5173
- [ ] Backend API responds at http://localhost:8080/api/v1/healthz
- [ ] Can register new user
- [ ] Can login with created user
- [ ] Dashboard shows sessions from database

## 🔒 Security Notes

- Passwords are hashed with bcrypt
- SQL injection prevented with prepared statements
- JWT tokens required for protected endpoints
- Role-based access control implemented
- Change JWT_SECRET in production
- Don't commit .env file with real credentials

## 📊 Key Files Modified

### Backend Database
- ✅ `go.mod` - PostgreSQL driver added
- ✅ `internal/db/*` - Database services (NEW)
- ✅ `internal/httpapi/auth.go` - DB integration
- ✅ `internal/httpapi/student.go` - DB integration
- ✅ `internal/httpapi/sessions.go` - DB integration
- ✅ `internal/httpapi/exams.go` - DB integration
- ✅ `internal/httpapi/admin.go` - DB integration

### Configuration
- ✅ `internal/config/config.go` - Enhanced
- ✅ `cmd/api/main.go` - Enhanced logging
- ✅ `Dockerfile` - Better build process

### Docker & Documentation
- ✅ `docker-compose.yml` - Full stack (NEW)
- ✅ `backend/schema.sql` - Database schema (NEW)
- ✅ `.env.example` - Configuration (NEW)
- ✅ Multiple documentation files (NEW)

## 🎯 Next Steps

1. **Test the API**
   - Run registration endpoint
   - Test login flow
   - Verify database persistence

2. **Customize**
   - Update JWT_SECRET
   - Modify environment variables
   - Adjust connection pool settings

3. **Deploy**
   - Set up AWS RDS PostgreSQL
   - Configure CloudFront CDN
   - Deploy via Docker/ECS

4. **Enhance**
   - Add more API endpoints
   - Implement caching
   - Add monitoring
   - Write tests

## 📞 Support

If something doesn't work:
1. Check [SETUP_GUIDE.md](SETUP_GUIDE.md) for troubleshooting
2. Review [README.md](README.md) for full documentation
3. Check database connection: `psql -U postgres`
4. View Docker logs: `docker-compose logs`

## 🎉 You're All Set!

Your Scalable Learning Platform with PostgreSQL is ready to use!

**Start now:**
```bash
docker-compose up
```

Then visit:
- **Frontend:** http://localhost:5173
- **Backend API:** http://localhost:8080/api/v1/healthz
- **Database:** localhost:5432

Happy coding! 🚀
