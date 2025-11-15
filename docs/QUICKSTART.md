# 快速入门指南

## 5分钟上手

### 步骤 1: 克隆项目

```bash
git clone <repository-url>
cd packet-capture-tool
```

### 步骤 2: 安装依赖

```bash
# Windows 命令行
install-deps.bat

# 或使用 PowerShell/命令行
go mod download
```

### 步骤 3: 运行程序

```bash
# 方式 1: 直接运行（开发模式）
run-dev.bat

# 方式 2: 编译后运行
build.bat
packet-capture.exe
```

### 步骤 4: 安装证书

程序首次启动时会提示安装根证书：
1. 点击 "是" 允许安装（需要管理员权限）
2. 或点击 GUI 中的 "重新安装证书" 按钮

### 步骤 5: 启用代理

**方法 A: 使用系统托盘**
1. 在任务栏找到程序图标（右下角）
2. 右键点击 → 选择 "✅ 启用代理"

**方法 B: 手动设置**
1. Windows 设置 → 网络和 Internet → 代理
2. 手动代理设置
3. 地址: `127.0.0.1`，端口: `8888`

### 步骤 6: 开始抓包

1. 打开浏览器（Chrome、Edge 等）
2. 访问任意网站，如: `https://github.com`
3. 在 Packet Capture Tool 界面中查看捕获的请求

### 步骤 7: 搜索和过滤

在搜索框输入关键词，例如：
- `github` - 查找 GitHub 相关请求
- `GET` - 查找所有 GET 请求
- `200` - 查找状态码为 200 的请求

### 步骤 8: 使用完毕后禁用代理

右键系统托盘图标 → "❌ 禁用代理"

## 常见操作

### 查看请求详情

点击请求列表中的任意行，可以看到：
- 完整 URL
- 请求头和响应头
- 请求体和响应体

### 清空请求列表

点击 "清空请求" 按钮

### 修改代理端口

1. 在 "代理端口" 输入框输入新端口号
2. 程序自动重启代理服务器
3. 重新启用系统代理

### 禁用 HTTPS 解密

取消勾选 "启用 HTTPS 解密"

## 测试示例

### 测试 HTTP 请求

```bash
# 确保代理已启用
curl -x http://127.0.0.1:8888 http://httpbin.org/get
```

### 测试 HTTPS 请求

```bash
# 需要先安装根证书
curl -x http://127.0.0.1:8888 https://httpbin.org/get
```

### 测试 POST 请求

```bash
curl -x http://127.0.0.1:8888 \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"test":"data"}' \
  https://httpbin.org/post
```

## 浏览器测试

### 测试 API 请求

1. 打开浏览器开发者工具（F12）
2. 在 Console 中执行：

```javascript
// 测试 GET 请求
fetch('https://api.github.com/users/github')
  .then(r => r.json())
  .then(data => console.log(data));

// 测试 POST 请求
fetch('https://httpbin.org/post', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ test: 'data' })
})
  .then(r => r.json())
  .then(data => console.log(data));
```

3. 在 Packet Capture Tool 中查看捕获的请求

## 进阶使用

### 抓取移动应用流量

**注意**: 这需要修改代码以监听所有网络接口

1. 修改 `internal/proxy/proxy.go`:
   ```go
   // 将 "127.0.0.1:port" 改为 "0.0.0.0:port"
   addr := fmt.Sprintf("0.0.0.0:%d", s.port)
   ```

2. 在防火墙中允许该端口

3. 在手机 WiFi 设置中配置代理:
   - 代理地址: 电脑的 IP 地址
   - 端口: 8888

4. 在手机上安装根证书:
   - 将 `%APPDATA%\PacketCaptureTool\root-ca.crt` 发送到手机
   - 在手机设置中安装证书

⚠️ **安全警告**: 监听 `0.0.0.0` 会将代理暴露到网络，请确保网络安全！

### 与 Postman 配合使用

1. 启用 Packet Capture Tool 代理
2. 打开 Postman
3. Settings → Proxy → 启用系统代理
4. 发送请求
5. 在 Packet Capture Tool 中查看完整的请求细节

### 调试 Node.js 应用

```bash
# 设置环境变量
set HTTP_PROXY=http://127.0.0.1:8888
set HTTPS_PROXY=http://127.0.0.1:8888

# 对于 HTTPS，还需要设置
set NODE_TLS_REJECT_UNAUTHORIZED=0

# 运行你的 Node.js 应用
node app.js
```

### 调试 Python 应用

```python
import os
import requests

# 设置代理
os.environ['HTTP_PROXY'] = 'http://127.0.0.1:8888'
os.environ['HTTPS_PROXY'] = 'http://127.0.0.1:8888'

# 禁用 SSL 验证（因为使用自签名证书）
response = requests.get('https://api.github.com', verify=False)
print(response.json())
```

## 常见问题

### Q: 为什么看不到 HTTPS 内容？

A: 
1. 检查证书是否安装: `certmgr.msc`
2. 确认 "启用 HTTPS 解密" 已勾选
3. 某些应用可能需要额外配置

### Q: 端口 8888 被占用怎么办？

A: 在 GUI 中修改端口号为其他值（如 8889、9999 等）

### Q: 浏览器显示证书错误

A: 
1. 重新安装证书
2. 重启浏览器
3. 清除浏览器缓存

### Q: Firefox 不走代理

A: Firefox 使用独立的代理设置，需要在 Firefox 设置中单独配置

### Q: 如何查看请求的请求体？

A: 点击请求列表中的行，在详情区域可以看到完整的请求体和响应体

## 性能建议

### 内存使用优化

程序默认保留 10000 条请求。如果内存不足：
- 定期点击 "清空请求"
- 使用搜索功能过滤无关请求

### 性能影响

- 启用代理会增加 ~10ms 延迟
- HTTPS 解密会增加额外的 TLS 握手时间
- 建议仅在需要调试时启用

## 下一步

- 阅读 [完整使用文档](USAGE.md)
- 了解 [系统架构](ARCHITECTURE.md)
- 查看 [项目 README](../README.md)

## 技术支持

遇到问题？
1. 查看 [FAQ 文档](USAGE.md#故障排除)
2. 提交 GitHub Issue
3. 查看日志输出寻找错误信息

## 反馈

如果这个工具对你有帮助，欢迎：
- ⭐ Star 项目
- 🐛 报告 Bug
- 💡 提出改进建议
- 🤝 贡献代码

---

Happy Debugging! 🎉
