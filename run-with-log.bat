@echo off
REM 运行程序并保存日志到文件

echo ========================================
echo  Packet Capture Tool - With Log File
echo ========================================
echo.

set CGO_ENABLED=1

echo Running program and saving logs to: app.log
echo.
echo Press Ctrl+C to stop the program
echo.

REM 运行程序，同时输出到终端和文件
go run main.go 2>&1 | tee app.log

echo.
echo Logs saved to: app.log
pause
