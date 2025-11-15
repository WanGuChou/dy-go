@echo off
REM 辅助脚本：切换到项目目录

echo ========================================
echo  Change to Project Directory
echo ========================================
echo.

REM 获取脚本所在目录
set "PROJECT_DIR=%~dp0"

echo Project directory: %PROJECT_DIR%
echo.

REM 切换到项目目录
cd /d "%PROJECT_DIR%"

echo Current directory: %CD%
echo.

REM 列出关键文件
echo Files in directory:
dir /b main.go go.mod build.bat 2>nul

if exist "main.go" (
    echo.
    echo ✓ You are now in the project directory
    echo.
    echo You can now run:
    echo   - run-debug.bat
    echo   - build.bat
    echo   - check-env.bat
) else (
    echo.
    echo ❌ Warning: main.go not found!
    echo This may not be the correct directory.
)

echo.
pause
