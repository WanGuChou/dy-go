@echo off
REM 构建 Packet Capture Tool

echo Building Packet Capture Tool...

REM 设置环境变量
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1

REM 编译
go build -ldflags="-H windowsgui" -o packet-capture.exe

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Build successful! Executable: packet-capture.exe
    echo.
    echo Usage:
    echo   1. Run packet-capture.exe
    echo   2. Install the root certificate when prompted
    echo   3. Use the system tray to enable/disable proxy
    echo.
) else (
    echo.
    echo Build failed! Please check the error messages above.
    echo.
)

pause
