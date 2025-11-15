# Windows 环境配置指南

## 问题诊断

如果你遇到了这个错误：
```
build constraints exclude all Go files in github.com\go-gl\gl@...
```

**原因**: Fyne GUI 框架需要 CGO 支持，而 CGO 需要 GCC 编译器。

## 解决方案

### 步骤 1: 安装 GCC 编译器

你需要安装一个 Windows 上的 GCC 编译器。有两个推荐选项：

#### 选项 A: TDM-GCC（推荐，最简单）

1. **下载 TDM-GCC**:
   - 访问: https://jmeubank.github.io/tdm-gcc/
   - 下载 `tdm64-gcc-x.x.x.exe`（64位版本）

2. **安装**:
   - 运行安装程序
   - 选择 "Create"（创建新安装）
   - 安装路径建议: `C:\TDM-GCC-64`
   - ✅ **重要**: 确保勾选 "Add to PATH"

3. **验证安装**:
   ```bash
   # 打开新的命令提示符窗口
   gcc --version
   ```
   
   应该显示类似：
   ```
   gcc.exe (tdm64-1) 10.3.0
   ```

#### 选项 B: MinGW-w64

1. **下载 MinGW-w64**:
   - 访问: https://www.mingw-w64.org/downloads/
   - 或使用在线安装器: https://sourceforge.net/projects/mingw-w64/

2. **安装**:
   - 运行安装程序
   - 选择 Architecture: `x86_64`
   - 选择 Threads: `posix`
   - 安装到: `C:\mingw-w64`

3. **添加到 PATH**:
   - 右键 "此电脑" → 属性 → 高级系统设置
   - 环境变量 → 系统变量 → Path
   - 添加: `C:\mingw-w64\mingw64\bin`

4. **验证**:
   ```bash
   gcc --version
   ```

### 步骤 2: 重新打开命令提示符

安装 GCC 后，**必须关闭并重新打开命令提示符**，以便 PATH 生效。

### 步骤 3: 运行环境检查

```bash
check-env.bat
```

这个脚本会检查所有必需的工具是否正确安装。

### 步骤 4: 安装 Go 依赖

```bash
install-deps.bat
```

### 步骤 5: 编译项目

```bash
build.bat
```

## 详细的错误排查

### 错误 1: gcc: command not found

**症状**:
```
gcc: command not found
```

**解决方案**:
1. 确保 GCC 已安装
2. 检查 GCC 是否在 PATH 中
3. 重新打开命令提示符

**验证**:
```bash
where gcc
# 应显示: C:\TDM-GCC-64\bin\gcc.exe
```

### 错误 2: cgo: C compiler "gcc" not found

**症状**:
```
cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in %PATH%
```

**解决方案**:
```bash
# 检查 PATH
echo %PATH%

# 手动添加 GCC 到 PATH
set PATH=%PATH%;C:\TDM-GCC-64\bin

# 或者永久添加（需要管理员权限）
setx PATH "%PATH%;C:\TDM-GCC-64\bin"
```

### 错误 3: build constraints exclude all Go files

**症状**:
```
build constraints exclude all Go files in github.com\go-gl\gl
```

**解决方案**:
```bash
# 确保 CGO 已启用
set CGO_ENABLED=1

# 清理缓存
go clean -cache
go clean -modcache

# 重新安装依赖
go mod download
go mod tidy

# 重新编译
go build
```

### 错误 4: undefined reference to `WinMain'

**症状**:
```
undefined reference to `WinMain'
```

**解决方案**:
这是正常的，使用以下命令之一：

```bash
# 开发模式（带控制台）
go build -o packet-capture.exe

# 发布模式（无控制台）
go build -ldflags="-H windowsgui" -o packet-capture.exe
```

## 完整的安装流程

### 从零开始（新电脑）

1. **安装 Go**:
   - https://golang.org/dl/
   - 下载并安装 Windows 64-bit 版本
   - 安装时选择 "Add to PATH"

2. **安装 GCC** (TDM-GCC):
   - https://jmeubank.github.io/tdm-gcc/
   - 下载并安装 64-bit 版本
   - 安装时选择 "Add to PATH"

3. **安装 Git** (可选但推荐):
   - https://git-scm.com/
   - 使用默认设置安装

4. **重启电脑** (或至少重新打开命令提示符)

5. **验证安装**:
   ```bash
   go version
   gcc --version
   git --version
   ```

6. **克隆项目** (或解压项目文件):
   ```bash
   git clone <repository-url>
   cd packet-capture-tool
   ```

7. **运行环境检查**:
   ```bash
   check-env.bat
   ```

8. **安装依赖**:
   ```bash
   install-deps.bat
   ```

9. **编译项目**:
   ```bash
   build.bat
   ```

10. **运行程序**:
    ```bash
    packet-capture-debug.exe
    ```

## 开发模式 vs 发布模式

### 开发模式（推荐调试时使用）

```bash
# 保留控制台窗口，可以看到日志
go build -o packet-capture-debug.exe
packet-capture-debug.exe
```

**优点**:
- 可以看到详细的日志输出
- 便于调试和排查问题
- 可以看到错误信息

### 发布模式（正式使用）

```bash
# 无控制台窗口，更专业
go build -ldflags="-H windowsgui" -o packet-capture.exe
packet-capture.exe
```

**优点**:
- 没有黑色的控制台窗口
- 看起来更专业
- 用户体验更好

## 性能优化建议

### 加快编译速度

```bash
# 使用并行编译
go build -p 4

# 启用编译缓存（默认启用）
go env GOCACHE

# 如果缓存损坏，清理缓存
go clean -cache
```

### 减小可执行文件大小

```bash
# 移除调试信息
go build -ldflags="-s -w" -o packet-capture.exe

# 进一步使用 UPX 压缩（需要单独安装 UPX）
upx --best packet-capture.exe
```

## 常见问题 FAQ

### Q: 为什么需要 GCC？

A: Fyne GUI 框架使用了 C 语言库（OpenGL），Go 通过 CGO 调用这些库。CGO 需要 C 编译器。

### Q: 可以使用 MSVC 或 Clang 吗？

A: 理论上可以，但 MinGW/TDM-GCC 是最简单和最兼容的选择。

### Q: 编译很慢，正常吗？

A: 首次编译会慢（5-10分钟），因为需要编译 CGO 依赖。后续编译会快很多（使用缓存）。

### Q: 可以交叉编译吗？

A: CGO 使得交叉编译变得复杂。建议在目标平台上直接编译。

### Q: 我的杀毒软件报警怎么办？

A: 这是误报。自己编译的程序可能被误判。添加到白名单即可。

## 还有问题？

1. 查看 [BUILD.md](BUILD.md) 了解更多构建选项
2. 查看 [QUICKSTART.md](QUICKSTART.md) 了解使用方法
3. 在 GitHub 上提交 Issue

## 参考链接

- Go 官网: https://golang.org/
- TDM-GCC: https://jmeubank.github.io/tdm-gcc/
- MinGW-w64: https://www.mingw-w64.org/
- Fyne 文档: https://developer.fyne.io/
- CGO 文档: https://golang.org/cmd/cgo/

---

祝编译顺利！如果这个指南帮到了你，请给项目点个 Star ⭐
