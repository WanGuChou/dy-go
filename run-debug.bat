@echo off
REM 调试模式运行，显示详细日志

echo ========================================
echo  Packet Capture Tool - Debug Mode
echo ========================================
echo.

REM 切换到脚本所在目录
cd /d "%~dp0"

REM 显示当前目录
echo Current directory: %CD%
echo.

REM 检查 main.go 是否存在
if not exist "main.go" (
    echo ❌ Error: main.go not found!
    echo.
    echo Current directory: %CD%
    echo.
    echo Please make sure you are in the project root directory.
    echo The directory should contain: main.go, go.mod, internal\, etc.
    echo.
    pause
    exit /b 1
)

echo ✓ Found main.go
echo.

REM 设置CGO环境变量
set CGO_ENABLED=1

REM 检查GCC
where gcc >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ❌ GCC not found! Please install TDM-GCC first.
    echo See: docs\WINDOWS_SETUP.md
    echo.
    pause
    exit /b 1
)

echo ✓ GCC found
echo.

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
