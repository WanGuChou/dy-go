@echo off
REM 安装项目依赖

echo Installing dependencies...
echo.

go mod download
go mod tidy

echo.
echo Dependencies installed successfully!
echo.

pause
