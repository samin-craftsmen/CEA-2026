@echo off
echo =====================================
echo   Starting Meal Planner Backend...
echo =====================================

echo.
echo Checking Go installation...
go version
IF %ERRORLEVEL% NEQ 0 (
    echo Go is not installed or not added to PATH.
    pause
    exit /b
)

echo.
echo Downloading dependencies...
go mod tidy

echo.
echo Starting server...
go run main.go

pause
