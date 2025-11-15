@echo off
echo ========================================
echo  GUI Test - Simple Window
echo ========================================
echo.
echo This will test if Fyne GUI works on your system.
echo A simple window should appear with "Hello! GUI works!"
echo.
echo Compiling test...

cd /d "%~dp0"
set CGO_ENABLED=1

REM 先编译
go build -o test-gui-app.exe test-gui.go

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ❌ Compilation failed!
    echo This means there's a problem with the build environment.
    pause
    exit /b 1
)

echo ✓ Compilation successful!
echo.
echo Running GUI test...
echo.
echo ======== LAUNCHING GUI ========
echo.
echo If a window appears, GUI works! ✓
echo If no window appears, there may be a graphics issue. ❌
echo.

test-gui-app.exe

echo.
echo ======== GUI TEST ENDED ========
echo.

REM 清理
if exist "test-gui-app.exe" del test-gui-app.exe

echo Did you see a window? (Y/N)
set /p answer="> "

if /i "%answer%"=="Y" (
    echo.
    echo ✓ Great! GUI works on your system.
    echo The main program should also work.
    echo Try running: run-simple.bat
) else (
    echo.
    echo ❌ GUI did not appear. Possible issues:
    echo   1. Graphics drivers
    echo   2. Remote desktop session
    echo   3. Virtual machine without GPU
    echo   4. Display server issues
)

echo.
pause
