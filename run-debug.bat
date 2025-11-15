@echo off
REM 调试模式运行，显示详细日志

echo ========================================
echo  Packet Capture Tool - Debug Mode
echo ========================================
echo.

REM 设置CGO环境变量
set CGO_ENABLED=1

REM 检查GCC
where gcc >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ❌ GCC not found! Please install TDM-GCC first.
    echo See: docs\WINDOWS_SETUP.md
    pause
    exit /b 1
)

echo Running in debug mode with detailed logs...
echo.
echo ----------------------------------------
echo Application Output:
echo ----------------------------------------
echo.

REM 运行程序，日志会显示在终端
go run main.go

echo.
echo ----------------------------------------
echo Program exited
echo ----------------------------------------
pause
