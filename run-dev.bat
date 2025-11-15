@echo off
REM 开发模式运行 Packet Capture Tool

echo Starting Packet Capture Tool in development mode...
echo.

REM 切换到脚本所在目录
cd /d "%~dp0"

REM 检查 main.go 是否存在
if not exist "main.go" (
    echo ❌ Error: main.go not found in current directory!
    echo Current directory: %CD%
    echo.
    pause
    exit /b 1
)

REM 设置CGO环境变量（Fyne需要）
set CGO_ENABLED=1

REM 检查GCC是否安装
where gcc >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ❌ 错误: 未找到 GCC 编译器！
    echo.
    echo Fyne GUI 框架需要 CGO 支持，必须安装 GCC。
    echo.
    echo 请按照以下步骤安装：
    echo.
    echo 方式 1 - TDM-GCC ^(推荐^):
    echo   1. 访问: https://jmeubank.github.io/tdm-gcc/
    echo   2. 下载并安装 64-bit 版本
    echo   3. 安装时确保选中 "Add to PATH"
    echo.
    echo 方式 2 - MinGW-w64:
    echo   1. 访问: https://www.mingw-w64.org/
    echo   2. 下载并安装
    echo   3. 将 bin 目录添加到系统 PATH
    echo.
    echo 安装完成后，重新打开命令提示符并运行此脚本。
    echo.
    pause
    exit /b 1
)

echo ✅ GCC 已安装
gcc --version | findstr "gcc"
echo.

echo 正在编译并运行...
go run main.go

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ❌ 运行失败！
    echo.
    echo 如果遇到问题，请查看: docs\BUILD.md
    echo.
)

pause
