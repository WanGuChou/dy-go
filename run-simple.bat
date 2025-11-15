@echo off
REM 最简化运行模式

echo ========================================
echo  Simple Run Mode
echo ========================================
echo.

cd /d "%~dp0"

if not exist "main.go" (
    echo Error: main.go not found!
    pause
    exit /b 1
)

set CGO_ENABLED=1

echo Step 1: Cleaning old build cache...
go clean -cache
echo.

echo Step 2: Building executable...
echo (This will take a few minutes on first run)
echo.
go build -o app-test.exe main.go

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ❌ Build failed! Check errors above.
    pause
    exit /b 1
)

echo.
echo ✓ Build successful!
echo.

echo Step 3: Running the program...
echo.
echo Look for:
echo   - A window titled "Packet Capture Tool"
echo   - An icon in system tray (bottom-right)
echo.
echo Starting in 3 seconds...
timeout /t 3 >nul

echo.
echo ======== PROGRAM IS STARTING ========
echo.

start "" app-test.exe

echo.
echo Program has been launched!
echo.
echo If you don't see a window:
echo   1. Check Task Manager (Ctrl+Shift+Esc) for app-test.exe
echo   2. Check system tray (bottom-right corner)
echo   3. Try Alt+Tab to find the window
echo.
echo Press any key to cleanup and exit...
pause >nul

echo.
echo Cleaning up...
if exist "app-test.exe" del app-test.exe

echo Done!
