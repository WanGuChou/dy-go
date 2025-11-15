@echo off
REM 构建 Packet Capture Tool

echo ========================================
echo  Packet Capture Tool - Build Script
echo ========================================
echo.

REM 切换到脚本所在目录
cd /d "%~dp0"

REM 检查 main.go 是否存在
if not exist "main.go" (
    echo ❌ Error: main.go not found!
    echo Current directory: %CD%
    echo.
    echo Please run this script from the project root directory.
    pause
    exit /b 1
)

echo Current directory: %CD%
echo.

REM 检查GCC是否安装
echo [1/4] Checking GCC...
where gcc >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ❌ 错误: 未找到 GCC 编译器！
    echo.
    echo Fyne GUI 框架需要 CGO 支持，必须安装 GCC。
    echo.
    echo 请安装 TDM-GCC 或 MinGW-w64
    echo 详细说明请查看: docs\BUILD.md
    echo.
    pause
    exit /b 1
)
echo ✅ GCC found
gcc --version | findstr "gcc"
echo.

REM 检查Go是否安装
echo [2/4] Checking Go...
where go >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ❌ 错误: 未找到 Go 编译器！
    echo 请从 https://golang.org/dl/ 下载并安装 Go
    echo.
    pause
    exit /b 1
)
echo ✅ Go found
go version
echo.

REM 设置环境变量
echo [3/4] Setting environment...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
echo ✅ CGO_ENABLED=1
echo.

REM 编译（带控制台窗口，方便调试）
echo [4/4] Building...
echo Building with console window (for debugging)...
go build -o packet-capture-debug.exe

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ❌ Build failed!
    echo.
    echo 常见问题:
    echo   1. GCC 未正确安装或不在 PATH 中
    echo   2. Go 模块依赖未下载（运行 install-deps.bat）
    echo   3. 缺少必要的系统库
    echo.
    echo 详细说明请查看: docs\BUILD.md
    echo.
    pause
    exit /b 1
)

echo ✅ Build successful: packet-capture-debug.exe
echo.

REM 编译无控制台版本
echo Building GUI version (without console)...
go build -ldflags="-H windowsgui" -o packet-capture.exe

if %ERRORLEVEL% EQU 0 (
    echo ✅ Build successful: packet-capture.exe
    echo.
    echo ========================================
    echo Build Complete!
    echo ========================================
    echo.
    echo 生成的文件:
    echo   - packet-capture-debug.exe  ^(带控制台，推荐调试用^)
    echo   - packet-capture.exe        ^(无控制台，正式使用^)
    echo.
    echo 使用方法:
    echo   1. 运行 packet-capture-debug.exe
    echo   2. 首次运行会提示安装证书（需要管理员权限）
    echo   3. 右键点击系统托盘图标启用代理
    echo   4. 开始抓包！
    echo.
    echo 详细使用说明: docs\QUICKSTART.md
    echo.
) else (
    echo ⚠️  GUI version build failed, but debug version is available
    echo.
)

pause
