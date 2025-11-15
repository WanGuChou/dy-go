@echo off
REM 捕获所有错误信息

echo ========================================
echo  Run and Catch All Errors
echo ========================================
echo.

cd /d "%~dp0"
set CGO_ENABLED=1

echo Building...
go build -o app-error-test.exe main.go

if %ERRORLEVEL% NEQ 0 (
    echo Build failed!
    pause
    exit /b 1
)

echo.
echo Running and logging all output to error-log.txt...
echo.

REM 运行并捕获所有输出
app-error-test.exe > error-log.txt 2>&1

echo.
echo Program terminated.
echo.
echo Checking error log...
echo.
echo ======== ERROR LOG CONTENT ========
type error-log.txt
echo ===================================
echo.

echo Error log saved to: error-log.txt
echo.

if exist "app-error-test.exe" del app-error-test.exe

pause
