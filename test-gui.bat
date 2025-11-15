@echo off
echo Testing GUI functionality...
echo.
set CGO_ENABLED=1
go run test-gui.go
pause
