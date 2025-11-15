@echo off
echo ========================================
echo Testing GCC in Current Environment
echo ========================================
echo.

echo Current PATH:
echo %PATH%
echo.

echo ----------------------------------------
echo.

echo Searching for GCC:
where gcc
echo.

echo ----------------------------------------
echo.

echo GCC Version:
gcc --version
echo.

echo ========================================
echo Test Complete
echo ========================================

pause
