@echo off
REM 检查程序是否正在运行

echo ========================================
echo  Check if Program is Running
echo ========================================
echo.

echo Checking for Go processes...
tasklist | findstr /i "go.exe main.exe packet-capture"

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ✓ Found running processes!
    echo.
    echo The program is likely running. Try:
    echo   1. Press Alt+Tab to find the GUI window
    echo   2. Check system tray ^(bottom-right corner^)
    echo.
) else (
    echo.
    echo ❌ No related processes found.
    echo The program is not running.
    echo.
)

echo.
echo Full process list:
tasklist | findstr /i "go main packet"

echo.
pause
