@echo off
REM Scalable Learning Platform - Quick Start Script for Windows

echo.
echo 🚀 Scalable Learning Platform - Quick Start
echo ===========================================
echo.

REM Check if Docker is installed
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker is not installed. Please install Docker first.
    pause
    exit /b 1
)

REM Check if Docker Compose is installed
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker Compose is not installed. Please install Docker Compose first.
    pause
    exit /b 1
)

echo ✅ Docker and Docker Compose are installed
echo.

REM Check if .env file exists
if not exist ".env" (
    echo 📝 Creating .env file from template...
    copy .env.example .env
    echo ✅ .env file created. Please review and update if needed.
)

echo.
echo 🐳 Starting Docker Compose services...
echo.

REM Start services
docker-compose up

echo.
echo ✅ Services are running!
echo.
echo 📍 Access the application at:
echo    Frontend: http://localhost:5173
echo    Backend:  http://localhost:8080
echo    Database: localhost:5432
echo.
echo To stop services, press Ctrl+C
pause
