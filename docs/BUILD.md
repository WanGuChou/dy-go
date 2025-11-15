# 构建指南

## 环境要求

### Windows 平台

**必需工具:**
- Go 1.21 或更高版本
- GCC 编译器（用于 CGO）
  - 推荐使用 [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)
  - 或者 [MinGW-w64](https://www.mingw-w64.org/)

**可选工具:**
- Git（用于版本控制）
- Make（可选，用于使用 Makefile）

### 安装 Go

1. 下载 Go 安装包: https://golang.org/dl/
2. 运行安装程序
3. 验证安装:
   ```bash
   go version
   ```

### 安装 GCC（CGO 支持）

Fyne GUI 需要 CGO，因此需要 GCC 编译器。

**方式 1: TDM-GCC（推荐）**

1. 下载 TDM-GCC: https://jmeubank.github.io/tdm-gcc/
2. 安装 64-bit 版本
3. 确保 GCC 在 PATH 中

**方式 2: MinGW-w64**

1. 下载 MinGW-w64
2. 安装并添加到 PATH
3. 验证安装:
   ```bash
   gcc --version
   ```

## 构建步骤

### 1. 克隆项目

```bash
git clone <repository-url>
cd packet-capture-tool
```

### 2. 安装依赖

```bash
# Windows
install-deps.bat

# 或手动执行
go mod download
go mod tidy
```

### 3. 编译项目

**方式 1: 使用批处理脚本**

```bash
build.bat
```

**方式 2: 使用 Make**

```bash
make build
```

**方式 3: 手动编译**

```bash
# 编译为控制台程序（开发模式）
go build -o packet-capture.exe

# 编译为 Windows GUI 程序（无控制台窗口）
go build -ldflags="-H windowsgui" -o packet-capture.exe
```

### 4. 运行程序

```bash
# 直接运行
packet-capture.exe

# 或开发模式（不编译）
go run main.go
```

## 构建选项

### 标准构建

```bash
go build -o packet-capture.exe
```

这会生成一个带控制台窗口的可执行文件，方便查看日志。

### 无控制台窗口

```bash
go build -ldflags="-H windowsgui" -o packet-capture.exe
```

适用于正式发布，没有黑色的控制台窗口。

### 带调试信息

```bash
go build -gcflags="all=-N -l" -o packet-capture-debug.exe
```

保留调试符号，方便使用 delve 等工具调试。

### 优化构建（减小体积）

```bash
go build -ldflags="-s -w" -o packet-capture.exe
```

- `-s`: 移除符号表
- `-w`: 移除 DWARF 调试信息

可以进一步使用 UPX 压缩:

```bash
upx --best packet-capture.exe
```

### 交叉编译

虽然项目主要针对 Windows，但也可以为其他平台编译（需要修改代码）:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o packet-capture-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o packet-capture-macos
```

**注意**: 跨平台编译需要修改系统特定的代码（证书安装、系统代理设置等）。

## 故障排除

### 问题 1: gcc: command not found

**原因**: 没有安装 GCC 编译器

**解决方案**: 安装 TDM-GCC 或 MinGW-w64

### 问题 2: undefined reference to `WinMain'

**原因**: 链接器找不到入口点

**解决方案**: 确保使用正确的编译命令，或添加 `-H windowsgui` 标志

### 问题 3: package fyne.io/fyne/v2: cannot find package

**原因**: 依赖未安装

**解决方案**:

```bash
go mod download
go mod tidy
```

### 问题 4: cgo: C compiler "gcc" not found

**原因**: GCC 不在 PATH 中

**解决方案**: 将 GCC 安装目录添加到 PATH 环境变量

### 问题 5: 编译很慢

**原因**: CGO 编译需要较长时间

**优化方案**:
- 使用 `go build -buildmode=pie` 增量编译
- 使用编译缓存（Go 1.10+ 默认启用）
- 升级到更快的 SSD

## 开发构建

### 快速迭代

在开发过程中，使用 `go run` 而不是每次都编译:

```bash
go run main.go
```

### 自动重新编译

使用 `air` 工具实现热重载:

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 创建 .air.toml 配置文件
air init

# 运行
air
```

### 启用 Race Detector

检测并发问题:

```bash
go run -race main.go
```

**注意**: Race detector 会显著降低性能，仅用于开发测试。

## 测试

### 运行单元测试

```bash
# 所有测试
go test ./...

# 特定包
go test ./internal/storage
go test ./internal/config

# 显示详细输出
go test -v ./...

# 测试覆盖率
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 运行集成测试

```bash
# 跳过需要代理服务器运行的测试
go test -short ./...
```

## 性能分析

### CPU Profile

```bash
go build -o packet-capture.exe
packet-capture.exe -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

### Memory Profile

```bash
go build -o packet-capture.exe
packet-capture.exe -memprofile=mem.prof
go tool pprof mem.prof
```

## 发布构建

### 创建发布版本

```bash
# 清理旧构建
del packet-capture.exe

# 编译优化版本
go build -ldflags="-s -w -H windowsgui" -o packet-capture.exe

# 压缩（可选）
upx --best packet-capture.exe

# 创建发布包
mkdir release
copy packet-capture.exe release\
copy README.md release\
copy docs\QUICKSTART.md release\
```

### 打包为安装程序

使用 [Inno Setup](https://jrsoftware.org/isinfo.php) 创建 Windows 安装程序。

示例脚本 (`setup.iss`):

```iss
[Setup]
AppName=Packet Capture Tool
AppVersion=1.0
DefaultDirName={pf}\PacketCaptureTool
DefaultGroupName=Packet Capture Tool
OutputDir=output
OutputBaseFilename=packet-capture-setup

[Files]
Source: "packet-capture.exe"; DestDir: "{app}"
Source: "README.md"; DestDir: "{app}"

[Icons]
Name: "{group}\Packet Capture Tool"; Filename: "{app}\packet-capture.exe"
Name: "{group}\Uninstall"; Filename: "{uninstallexe}"
```

## CI/CD 集成

### GitHub Actions

创建 `.github/workflows/build.yml`:

```yaml
name: Build

on: [push, pull_request]

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install dependencies
        run: |
          go mod download
          go mod tidy
      
      - name: Run tests
        run: go test -v ./...
      
      - name: Build
        run: go build -v -o packet-capture.exe
      
      - name: Upload artifact
        uses: actions/upload-artifact@v3
        with:
          name: packet-capture-windows
          path: packet-capture.exe
```

## 依赖管理

### 更新依赖

```bash
# 查看可更新的依赖
go list -u -m all

# 更新所有依赖
go get -u ./...
go mod tidy

# 更新特定依赖
go get -u fyne.io/fyne/v2
go mod tidy
```

### 查看依赖树

```bash
go mod graph
```

### 清理未使用的依赖

```bash
go mod tidy
```

## 常见编译问题

### Fyne 相关

**问题**: Fyne 编译失败

**解决方案**: 
1. 确保 GCC 已安装
2. 在 Windows 上使用 TDM-GCC
3. 检查 PATH 环境变量

### CGO 相关

**问题**: CGO_ENABLED=1 required

**解决方案**:

```bash
set CGO_ENABLED=1
go build
```

### 内存不足

**问题**: 编译时内存不足

**解决方案**:
```bash
# 限制并行编译数
go build -p 2
```

## 调试构建

### 使用 Delve 调试

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试运行
dlv debug main.go

# 常用命令
(dlv) break main.main
(dlv) continue
(dlv) next
(dlv) print var
```

### Visual Studio Code

创建 `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/main.go"
    }
  ]
}
```

## 构建优化建议

1. **使用最新版本的 Go**: 新版本通常有更好的性能
2. **启用编译缓存**: Go 1.10+ 默认启用
3. **使用并行编译**: `-p` 标志控制并行度
4. **减少依赖**: 只引入必需的包
5. **静态链接**: 减少运行时依赖

## 更多资源

- [Go 官方文档](https://golang.org/doc/)
- [Fyne 文档](https://developer.fyne.io/)
- [CGO 文档](https://golang.org/cmd/cgo/)
- [Go 构建模式](https://pkg.go.dev/cmd/go#hdr-Build_modes)
