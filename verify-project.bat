@echo off
REM 验证项目完整性

echo ========================================
echo Packet Capture Tool - Project Verification
echo ========================================
echo.

echo [1/6] Checking Go installation...
go version >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Go is not installed or not in PATH
    exit /b 1
) else (
    go version
    echo ✅ Go is installed
)
echo.

echo [2/6] Checking project structure...
if not exist "main.go" (
    echo ❌ main.go not found
    exit /b 1
)
if not exist "internal\" (
    echo ❌ internal directory not found
    exit /b 1
)
echo ✅ Project structure looks good
echo.

echo [3/6] Checking Go modules...
if not exist "go.mod" (
    echo ❌ go.mod not found
    exit /b 1
)
if not exist "go.sum" (
    echo ⚠️  go.sum not found, will be created on first build
)
echo ✅ Go modules configured
echo.

echo [4/6] Verifying dependencies...
go mod verify
if %ERRORLEVEL% NEQ 0 (
    echo ⚠️  Some dependencies need to be downloaded
    echo Run: install-deps.bat
) else (
    echo ✅ Dependencies verified
)
echo.

echo [5/6] Running tests...
go test ./internal/config ./internal/storage
if %ERRORLEVEL% NEQ 0 (
    echo ⚠️  Some tests failed
) else (
    echo ✅ Tests passed
)
echo.

echo [6/6] Checking documentation...
set doc_count=0
if exist "README.md" set /a doc_count+=1
if exist "docs\ARCHITECTURE.md" set /a doc_count+=1
if exist "docs\BUILD.md" set /a doc_count+=1
if exist "docs\USAGE.md" set /a doc_count+=1
if exist "docs\QUICKSTART.md" set /a doc_count+=1
if exist "CONTRIBUTING.md" set /a doc_count+=1
if exist "PROJECT_SUMMARY.md" set /a doc_count+=1

echo ✅ Found %doc_count%/7 documentation files
echo.

echo ========================================
echo Verification Complete!
echo ========================================
echo.
echo Next steps:
echo   1. Run: install-deps.bat (if needed)
echo   2. Run: build.bat (to compile)
echo   3. Run: packet-capture.exe (to start)
echo.
echo For more information, see:
echo   - README.md
echo   - docs\QUICKSTART.md
echo.

pause
