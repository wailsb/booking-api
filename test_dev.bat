@echo off
setlocal enabledelayedexpansion

echo ===================================
echo Resetting and Starting Docker Stack...
echo ===================================
docker compose down -v
docker compose up --build -d

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [ERROR] Docker command failed. Ensure Docker Desktop is running!
    pause
    exit /b %ERRORLEVEL%
)

echo.
echo Waiting 6 seconds for Postgres and Migrations...
timeout /t 6 /nobreak > nul

set BASE_URL=http://localhost:8080/api/v1
set BOOKABLE_ID=a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11

echo.
echo ===================================
echo 1. GET Seeded Bookings for Bookable
echo ===================================
curl -X GET "%BASE_URL%/bookables/%BOOKABLE_ID%/bookings" -H "Accept: application/json"
echo.

echo.
echo ===================================
echo 2. TEST Overlapping Booking
echo ===================================
curl -X POST "%BASE_URL%/bookings" ^
  -H "Content-Type: application/json" ^
  -d "{\"bookable_id\":\"%BOOKABLE_ID%\",\"customer_name\":\"Overlap Test\",\"customer_email\":\"test@test.com\",\"start_time\":\"2026-10-01T09:30:00Z\",\"end_time\":\"2026-10-01T10:30:00Z\"}"
echo.

echo.
echo ===================================
echo Seed and Environment Ready
echo ===================================
pause