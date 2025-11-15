@echo off
REM 仅编译，不运行

echo ========================================
echo  Compile Only (with verbose output)
echo ========================================
echo.

cd /d "%~dp0"

if not exist "main.go" (
    echo Error: main.go not found!
    pause
    exit /b 1
)

set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

echo Compiling with verbose output...
echo This will show you what's happening during compilation.
echo.
echo ----------------------------------------
echo.

REM 详细编译过程
go build -v -x -o packet-capture-test.exe main.go

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ========================================
    echo ❌ Compilation FAILED!
    echo ========================================
    echo.
    echo Please check the error messages above.
    echo.
) else (
    echo.
    echo ========================================
    echo ✓ Compilation SUCCESSFUL!
    echo ========================================
    echo.
    echo Generated file: packet-capture-test.exe
    echo File size:
    dir packet-capture-test.exe | findstr "packet-capture-test.exe"
    echo.
    echo You can now run: packet-capture-test.exe
    echo.
)

pause
