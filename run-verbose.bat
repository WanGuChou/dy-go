@echo off
REM 超详细模式运行

echo ========================================
echo  Packet Capture Tool - VERBOSE Mode
echo ========================================
echo.

cd /d "%~dp0"

if not exist "main.go" (
    echo Error: main.go not found!
    pause
    exit /b 1
)

set CGO_ENABLED=1

echo Step 1: Compiling (this may take a few minutes)...
echo.

REM 先编译，显示编译过程
go build -v -o packet-capture-temp.exe main.go

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ❌ Compilation failed!
    pause
    exit /b 1
)

echo.
echo ✓ Compilation successful!
echo.
echo Step 2: Running the program...
echo.
echo ----------------------------------------
echo Program Output:
echo ----------------------------------------
echo.

REM 运行编译好的程序
packet-capture-temp.exe

echo.
echo ----------------------------------------
echo Program terminated
echo ----------------------------------------
echo.

REM 清理临时文件
if exist "packet-capture-temp.exe" del packet-capture-temp.exe

pause
