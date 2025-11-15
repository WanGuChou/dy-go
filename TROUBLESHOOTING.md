# 常见问题排查指南

本文档汇总了使用 Packet Capture Tool 时可能遇到的常见问题及解决方案。

## 编译问题

### 问题 1: build constraints exclude all Go files

**完整错误信息**:
```
package command-line-arguments
        imports packet-capture-tool/internal/gui
        imports fyne.io/fyne/v2/app
        imports github.com/go-gl/gl/v2.1/gl: 
        build constraints exclude all Go files in C:\Users\...\github.com\go-gl\gl@...\v2.1\gl
```

**根本原因**: 
- Fyne GUI 框架需要 CGO 支持
- CGO 需要 C 编译器（GCC）
- 系统中没有安装 GCC 或 GCC 不在 PATH 中

**解决方案**:

**步骤 1: 安装 GCC**

选择以下任一方式：

**方式 A - TDM-GCC（推荐）**:
1. 访问: https://jmeubank.github.io/tdm-gcc/
2. 下载 `tdm64-gcc-x.x.x.exe`
3. 运行安装程序
4. ✅ 确保勾选 "Add to PATH"
5. 安装完成

**方式 B - MinGW-w64**:
1. 访问: https://www.mingw-w64.org/downloads/
2. 下载并安装
3. 手动添加到 PATH: `C:\mingw-w64\mingw64\bin`

**步骤 2: 验证安装**

关闭并重新打开命令提示符，然后运行：
```bash
gcc --version
```

应该显示类似：
```
gcc.exe (tdm64-1) 10.3.0
```

**步骤 3: 重新编译**

```bash
# 运行环境检查
check-env.bat

# 清理缓存
go clean -cache
go clean -modcache

# 重新安装依赖
install-deps.bat

# 编译
build.bat
```

**详细指南**: 查看 [WINDOWS_SETUP.md](docs/WINDOWS_SETUP.md)

---

### 问题 2: GCC 在 Cursor/VS Code 终端中无法识别

**症状**:
- 在 Windows 自带终端（cmd/PowerShell）中 `gcc --version` 正常
- 在 Cursor/VS Code 集成终端中显示：
  ```
  'gcc' is not recognized as an internal or external command
  ```

**原因**: 
- Cursor/VS Code 在启动时读取环境变量
- 如果在 IDE 启动后才安装 GCC，IDE 的终端不会自动更新 PATH

**解决方案**:

**方案 A - 重启 Cursor（推荐）**:
1. 完全关闭 Cursor（不是只关闭窗口）
2. 检查任务管理器，确保没有 Cursor 进程
3. 重新打开 Cursor
4. 在终端中测试：`gcc --version`

**方案 B - 临时设置 PATH**:
```bash
# 在 Cursor 终端中执行（根据实际路径修改）
set PATH=%PATH%;C:\TDM-GCC-64\bin

# 验证
gcc --version
```

**方案 C - 使用 Windows 终端**:
```bash
# 打开 Windows 命令提示符
cd C:\path\to\packet-capture-tool
build.bat
```

**验证环境变量**:

运行测试脚本：
```bash
test-gcc.bat
```

或手动检查：
```bash
# 查看 PATH
echo %PATH%

# 查找 GCC
where gcc

# 对比 Windows 终端和 Cursor 终端的输出
```

---

### 问题 3: cgo: C compiler "gcc" not found

**错误信息**:
```
cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in %PATH%
```

**原因**: GCC 没有添加到 PATH 环境变量

**解决方案**:

**临时方案**:
```bash
# 手动添加到当前会话
set PATH=%PATH%;C:\TDM-GCC-64\bin
```

**永久方案**:
1. 右键 "此电脑" → 属性
2. 高级系统设置 → 环境变量
3. 系统变量 → Path → 编辑
4. 添加: `C:\TDM-GCC-64\bin`（根据实际安装路径修改）
5. 确定 → 重新打开命令提示符

**验证**:
```bash
where gcc
# 应显示: C:\TDM-GCC-64\bin\gcc.exe
```

---

### 问题 4: undefined reference to WinMain

**错误信息**:
```
undefined reference to `WinMain'
```

**原因**: 这是 Windows GUI 程序的链接错误

**解决方案**:

对于开发调试，使用带控制台的版本：
```bash
go build -o packet-capture-debug.exe
```

对于发布版本，使用无控制台的版本：
```bash
go build -ldflags="-H windowsgui" -o packet-capture.exe
```

**推荐**: 直接使用 `build.bat`，它会同时生成两个版本。

---

### 问题 5: 编译速度很慢

**现象**: 首次编译需要 5-10 分钟

**原因**: CGO 需要编译 C 依赖库

**解决方案**:

这是正常的！首次编译会慢，后续编译会使用缓存，速度会快很多。

**优化建议**:
```bash
# 使用并行编译
go build -p 4

# 检查缓存位置
go env GOCACHE

# 如果缓存损坏，清理后重建
go clean -cache
go build
```

---

## 运行问题

### 问题 6: 证书安装失败

**错误信息**:
```
Failed to install certificate: exit status 1
```

**原因**: 需要管理员权限安装证书到系统信任存储

**解决方案**:

**方式 1**: 以管理员身份运行程序
- 右键 `packet-capture-debug.exe`
- 选择 "以管理员身份运行"

**方式 2**: 手动安装证书
```bash
# 以管理员身份打开命令提示符
cd %APPDATA%\PacketCaptureTool
certutil -addstore Root root-ca.crt
```

**验证安装**:
```bash
certutil -verifystore Root "Packet Capture Tool Root CA"
```

---

### 问题 7: 无法捕获 HTTPS 流量

**现象**: HTTP 请求可以看到，但 HTTPS 请求为空或失败

**可能原因**:

1. **证书未安装**
   ```bash
   # 检查证书
   certmgr.msc
   # 在 "受信任的根证书颁发机构" 中查找 "Packet Capture Tool Root CA"
   ```

2. **HTTPS 解密未启用**
   - 检查 GUI 中的 "启用 HTTPS 解密" 是否勾选

3. **浏览器使用自己的证书存储**
   - Chrome/Edge: 使用系统证书（自动生效）
   - Firefox: 使用独立证书存储（需要手动导入）

**Firefox 解决方案**:
1. 访问: `about:preferences#privacy`
2. 证书 → 查看证书 → 证书颁发机构
3. 导入 → 选择 `%APPDATA%\PacketCaptureTool\root-ca.crt`
4. 勾选 "信任此 CA 用于标识网站"

---

### 问题 8: 代理端口被占用

**错误信息**:
```
Failed to listen on 127.0.0.1:8888: bind: address already in use
```

**原因**: 端口 8888 被其他程序占用

**查找占用程序**:
```bash
netstat -ano | findstr :8888
```

**解决方案**:

**方式 1**: 关闭占用端口的程序
```bash
# 找到 PID（上一步输出的最后一列）
taskkill /PID <PID> /F
```

**方式 2**: 修改代理端口
1. 在 GUI 的 "代理端口" 输入框输入新端口（如 8889）
2. 程序会自动重启代理服务器
3. 更新系统代理设置中的端口号

---

### 问题 9: 浏览器显示证书错误

**现象**: 浏览器提示 "您的连接不是私密连接" 或类似警告

**可能原因**:

1. **证书未安装或安装不正确**
2. **浏览器缓存了旧的证书**
3. **证书已过期**

**解决方案**:

**步骤 1**: 重新安装证书
- 在 GUI 中点击 "重新安装证书"
- 或运行 `certutil -addstore Root %APPDATA%\PacketCaptureTool\root-ca.crt`

**步骤 2**: 清除浏览器缓存
- Chrome: `chrome://settings/clearBrowserData` → 选择 "缓存的图片和文件"
- Edge: `edge://settings/clearBrowserData`
- Firefox: `about:preferences#privacy` → 清除数据

**步骤 3**: 重启浏览器

---

### 问题 10: 系统代理设置不生效

**现象**: 启用代理后，浏览器流量没有经过代理

**检查步骤**:

**1. 验证代理设置**
```bash
# 查看注册表
reg query "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyEnable
reg query "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyServer
```

**2. 手动设置代理**
- Windows 设置 → 网络和 Internet → 代理
- 手动代理设置
- 地址: `127.0.0.1`，端口: `8888`

**3. Firefox 特殊情况**

Firefox 不使用系统代理，需要单独配置：
1. 设置 → 常规 → 网络设置
2. 手动代理配置
3. HTTP 代理: `127.0.0.1`，端口: `8888`
4. ✅ 勾选 "也将此代理用于 HTTPS"

---

### 问题 11: 程序崩溃或无响应

**现象**: 程序运行一段时间后崩溃或卡死

**可能原因**:

1. **内存占用过高**
   - 默认保留 10000 条请求
   - 定期点击 "清空请求"

2. **死锁或并发问题**
   - 记录完整的日志
   - 提交 Issue

**调试步骤**:

1. **使用调试版本**
   ```bash
   packet-capture-debug.exe > log.txt 2>&1
   ```

2. **检查日志**
   - 查看 `log.txt` 中的错误信息

3. **减少负载**
   - 使用搜索功能过滤请求
   - 定期清空请求列表
   - 避免同时打开大量网页

---

## 性能问题

### 问题 12: 网络速度变慢

**原因**: 代理会增加一定延迟

**优化建议**:

1. **仅在需要时启用代理**
   - 调试完成后立即禁用

2. **使用搜索过滤**
   - 减少 GUI 刷新压力

3. **禁用不需要的功能**
   - 如果只需要 HTTP，可以禁用 HTTPS 解密

4. **调整刷新频率**
   - 修改 `internal/gui/gui.go` 中的刷新间隔

---

### 问题 13: 内存占用过高

**现象**: 程序占用几百 MB 内存

**原因**: 捕获了大量请求（特别是包含大文件的请求）

**解决方案**:

1. **定期清空**
   - 点击 "清空请求" 按钮

2. **使用搜索过滤**
   - 只显示需要的请求

3. **调整内存限制**
   - 修改 `internal/storage/storage.go`:
     ```go
     // 将 10000 改为更小的值
     if len(s.requests) > 5000 {
         s.requests = s.requests[len(s.requests)-5000:]
     }
     ```

---

## 其他问题

### 问题 14: 杀毒软件报警

**现象**: Windows Defender 或其他杀毒软件提示威胁

**原因**: 自己编译的程序可能被误判为可疑软件

**解决方案**:

1. **添加到白名单**
   - Windows Defender: 设置 → 病毒和威胁防护 → 排除项
   - 添加 `packet-capture.exe` 所在目录

2. **验证程序安全性**
   - 你自己编译的程序是安全的
   - 查看源代码确认没有恶意行为

---

### 问题 15: 某些应用程序的流量抓不到

**原因**: 某些应用不使用系统代理

**常见应用**:

1. **Electron 应用** (如 VS Code, Slack)
   - 启动时添加参数: `--proxy-server=127.0.0.1:8888`

2. **Java 应用**
   ```bash
   java -Dhttp.proxyHost=127.0.0.1 -Dhttp.proxyPort=8888 -Dhttps.proxyHost=127.0.0.1 -Dhttps.proxyPort=8888 -jar app.jar
   ```

3. **Python 应用**
   ```python
   import os
   os.environ['HTTP_PROXY'] = 'http://127.0.0.1:8888'
   os.environ['HTTPS_PROXY'] = 'http://127.0.0.1:8888'
   ```

4. **Node.js 应用**
   ```bash
   set HTTP_PROXY=http://127.0.0.1:8888
   set HTTPS_PROXY=http://127.0.0.1:8888
   node app.js
   ```

---

## 获取帮助

如果以上方案都无法解决你的问题：

1. **查看完整文档**
   - [WINDOWS_SETUP.md](docs/WINDOWS_SETUP.md) - Windows 环境配置
   - [BUILD.md](docs/BUILD.md) - 构建指南
   - [USAGE.md](docs/USAGE.md) - 使用手册

2. **提交 Issue**
   - 访问 GitHub Issues
   - 提供以下信息：
     - 完整的错误信息
     - 操作系统版本
     - Go 版本 (`go version`)
     - GCC 版本 (`gcc --version`)
     - 详细的重现步骤

3. **查看日志**
   ```bash
   packet-capture-debug.exe > log.txt 2>&1
   ```
   将 `log.txt` 内容附在 Issue 中

---

## 快速诊断清单

遇到问题时，按照以下清单依次检查：

- [ ] Go 已安装且版本 >= 1.21
- [ ] GCC 已安装（TDM-GCC 或 MinGW-w64）
- [ ] GCC 在 PATH 中（`where gcc` 有输出）
- [ ] 运行过 `install-deps.bat`
- [ ] 证书已安装（`certmgr.msc` 可查看）
- [ ] 代理已启用（系统设置或托盘菜单）
- [ ] 端口没有被占用（`netstat -ano | findstr :8888`）
- [ ] 浏览器使用系统代理或已手动配置

如果以上全部通过，但仍有问题，请提交 Issue。

---

祝使用愉快！🎉
