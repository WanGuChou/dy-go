@echo off
REM 完整诊断工具

echo ========================================
echo  System Diagnostics for Packet Capture Tool
echo ========================================
echo.

cd /d "%~dp0"

echo [1/8] Checking files...
if exist "main.go" (echo ✓ main.go exists) else (echo ❌ main.go missing)
if exist "go.mod" (echo ✓ go.mod exists) else (echo ❌ go.mod missing)
if exist "internal\" (echo ✓ internal\ exists) else (echo ❌ internal\ missing)
echo.

echo [2/8] Checking Go...
go version
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Go not found!
    pause
    exit /b 1
)
echo.

echo [3/8] Checking GCC...
gcc --version | findstr "gcc"
if %ERRORLEVEL% NEQ 0 (
    echo ❌ GCC not found!
    pause
    exit /b 1
)
echo.

echo [4/8] Checking CGO...
set CGO_ENABLED=1
go env CGO_ENABLED
echo.

echo [5/8] Checking running processes...
tasklist | findstr /i "go.exe main.exe packet" >nul
if %ERRORLEVEL% EQU 0 (
    echo ⚠️  Found existing processes!
    tasklist | findstr /i "go.exe main.exe packet"
    echo.
    echo Please close them before running the program.
    echo.
) else (
    echo ✓ No conflicting processes found
)
echo.

echo [6/8] Checking display/graphics...
echo Display: %SESSIONNAME%
echo Graphics drivers: 
wmic path win32_VideoController get name
echo.

echo [7/8] Testing simple Go program...
echo package main > test-simple.go
echo import "fmt" >> test-simple.go
echo func main() { fmt.Println("Go works!") } >> test-simple.go
go run test-simple.go
del test-simple.go
echo.

echo [8/8] Testing Fyne GUI availability...
go list -m fyne.io/fyne/v2
if %ERRORLEVEL% NEQ 0 (
    echo ⚠️  Fyne module not found in go.mod
    echo Running: go mod download
    go mod download
)
echo.

echo ========================================
echo Diagnosis Complete
echo ========================================
echo.
echo Next steps:
echo   1. Run: test-gui.bat (test if GUI works at all)
echo   2. Run: run-verbose.bat (compile and run with full output)
echo.
pause
