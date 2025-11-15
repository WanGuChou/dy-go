@echo off
REM 检查开发环境

echo ========================================
echo  Environment Check for Packet Capture Tool
echo ========================================
echo.

set ERROR_COUNT=0

REM 检查 Go
echo [1/3] Checking Go installation...
where go >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Go is NOT installed
    echo    Download from: https://golang.org/dl/
    set /a ERROR_COUNT+=1
) else (
    echo ✅ Go is installed
    go version
)
echo.

REM 检查 GCC
echo [2/3] Checking GCC ^(required for Fyne CGO^)...
where gcc >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ❌ GCC is NOT installed
    echo    This is REQUIRED for Fyne GUI framework!
    echo.
    echo    Install options:
    echo    1. TDM-GCC ^(Recommended^): https://jmeubank.github.io/tdm-gcc/
    echo    2. MinGW-w64: https://www.mingw-w64.org/
    echo.
    set /a ERROR_COUNT+=1
) else (
    echo ✅ GCC is installed
    gcc --version | findstr "gcc"
)
echo.

REM 检查 Git (可选)
echo [3/3] Checking Git ^(optional^)...
where git >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ⚠️  Git is NOT installed ^(optional, but recommended^)
    echo    Download from: https://git-scm.com/
) else (
    echo ✅ Git is installed
    git --version
)
echo.

echo ========================================
echo Environment Check Complete
echo ========================================
echo.

if %ERROR_COUNT% EQU 0 (
    echo ✅ All required tools are installed!
    echo.
    echo Next steps:
    echo   1. Run: install-deps.bat  ^(install Go dependencies^)
    echo   2. Run: build.bat         ^(compile the project^)
    echo   3. Run: packet-capture-debug.exe
    echo.
) else (
    echo ❌ Found %ERROR_COUNT% missing tool^(s^)
    echo.
    echo Please install the missing tools and run this check again.
    echo.
    echo For detailed instructions, see: docs\BUILD.md
    echo.
)

pause
