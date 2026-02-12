#!/bin/bash

# Database migration and management script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Load environment variables
if [ -f ".env" ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

# Set defaults
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-scalable_learning}

# PGPASSWORD for non-interactive password entry
export PGPASSWORD=$DB_PASSWORD

echo -e "${GREEN}Database Management Script${NC}"
echo "=============================="
echo ""

# Parse command
COMMAND=${1:-help}

case $COMMAND in
    init)
        echo -e "${YELLOW}Initializing database...${NC}"
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f backend/schema.sql
        echo -e "${GREEN}✅ Database initialized successfully!${NC}"
        ;;
    
    create)
        echo -e "${YELLOW}Creating database...${NC}"
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME;"
        echo -e "${GREEN}✅ Database created successfully!${NC}"
        ;;
    
    drop)
        echo -e "${RED}WARNING: This will drop the database $DB_NAME${NC}"
        read -p "Are you sure? (y/N) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "DROP DATABASE IF EXISTS $DB_NAME;"
            echo -e "${GREEN}✅ Database dropped!${NC}"
        fi
        ;;
    
    reset)
        echo -e "${YELLOW}Resetting database...${NC}"
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "DROP DATABASE IF EXISTS $DB_NAME;"
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME;"
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f backend/schema.sql
        echo -e "${GREEN}✅ Database reset successfully!${NC}"
        ;;
    
    backup)
        BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
        echo -e "${YELLOW}Backing up database to $BACKUP_FILE...${NC}"
        pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME > $BACKUP_FILE
        echo -e "${GREEN}✅ Backup created: $BACKUP_FILE${NC}"
        ;;
    
    restore)
        if [ -z "$2" ]; then
            echo -e "${RED}Error: Please provide backup file${NC}"
            echo "Usage: ./db.sh restore <backup_file.sql>"
            exit 1
        fi
        echo -e "${YELLOW}Restoring database from $2...${NC}"
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$2"
        echo -e "${GREEN}✅ Database restored successfully!${NC}"
        ;;
    
    status)
        echo -e "${YELLOW}Database status:${NC}"
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
            SELECT 
                schemaname,
                tablename,
                pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
            FROM pg_tables
            WHERE schemaname = 'public'
            ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;"
        ;;
    
    connect)
        echo -e "${YELLOW}Connecting to database...${NC}"
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME
        ;;
    
    seed)
        echo -e "${YELLOW}Seeding sample data...${NC}"
        # Create sample data
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME <<EOF
-- Insert sample users
INSERT INTO users (name, email, password_hash, role) VALUES
('Admin User', 'admin@example.com', '\$2a\$10\$YourHashedPassword', 'ADMIN'),
('Instructor User', 'instructor@example.com', '\$2a\$10\$YourHashedPassword', 'INSTRUCTOR'),
('Student User 1', 'student1@example.com', '\$2a\$10\$YourHashedPassword', 'STUDENT'),
('Student User 2', 'student2@example.com', '\$2a\$10\$YourHashedPassword', 'STUDENT');

-- Insert sample sessions
INSERT INTO sessions (title, description, start_time, end_time, status, video_url, created_by) VALUES
('Introduction to Data Science', 'Basics of data science and machine learning', NOW() + INTERVAL '2 days', NOW() + INTERVAL '2 days 2 hours', 'SCHEDULED', 'https://example.com/video1', 2),
('Advanced SQL', 'Complex queries and optimization', NOW() - INTERVAL '1 days', NOW() - INTERVAL '23 hours', 'COMPLETED', 'https://example.com/video2', 2);

-- Insert sample exams
INSERT INTO exams (session_id, title, open_time, close_time, max_score, type) VALUES
(1, 'Data Science Quiz 1', NOW() + INTERVAL '3 days', NOW() + INTERVAL '4 days', 100, 'EXAM'),
(2, 'SQL Assignment 1', NOW() - INTERVAL '1 days', NOW() + INTERVAL '5 days', 50, 'ASSIGNMENT');

EOF
        echo -e "${GREEN}✅ Sample data seeded!${NC}"
        ;;
    
    *)
        echo "Usage: $0 {command} [options]"
        echo ""
        echo "Commands:"
        echo "  init      - Initialize database schema"
        echo "  create    - Create new database"
        echo "  drop      - Drop database"
        echo "  reset     - Reset database (drop and recreate)"
        echo "  backup    - Backup database"
        echo "  restore   - Restore from backup (requires file path)"
        echo "  status    - Show database status"
        echo "  connect   - Connect to database with psql"
        echo "  seed      - Seed sample data"
        echo ""
        echo "Environment variables (.env):"
        echo "  DB_HOST     - Database host (default: localhost)"
        echo "  DB_PORT     - Database port (default: 5432)"
        echo "  DB_USER     - Database user (default: postgres)"
        echo "  DB_PASSWORD - Database password (default: postgres)"
        echo "  DB_NAME     - Database name (default: scalable_learning)"
        echo ""
        ;;
esac
