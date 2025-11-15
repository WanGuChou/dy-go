# Makefile for Packet Capture Tool

.PHONY: all build run clean deps help

# 默认目标
all: build

# 构建项目
build:
	@echo "Building Packet Capture Tool..."
	@go build -ldflags="-H windowsgui" -o packet-capture.exe
	@echo "Build complete: packet-capture.exe"

# 开发模式运行
run:
	@echo "Running in development mode..."
	@go run main.go

# 安装依赖
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# 清理构建产物
clean:
	@echo "Cleaning build artifacts..."
	@if exist packet-capture.exe del packet-capture.exe
	@echo "Clean complete"

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

# 运行测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 显示帮助信息
help:
	@echo Packet Capture Tool - Makefile Help
	@echo.
	@echo Available targets:
	@echo   make build    - Build the executable
	@echo   make run      - Run in development mode
	@echo   make deps     - Install dependencies
	@echo   make clean    - Clean build artifacts
	@echo   make fmt      - Format code
	@echo   make test     - Run tests
	@echo   make help     - Show this help message
