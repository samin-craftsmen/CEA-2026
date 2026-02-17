@echo off
echo Installing frontend dependencies...
call npm install

echo Starting development server...
call npm run dev

pause
