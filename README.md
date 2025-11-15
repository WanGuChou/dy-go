# Packet Capture Tool - 本地抓包工具 🚀

一个轻量级的本地网络抓包工具，使用 Golang 开发，专为 Windows 平台设计，功能对标 Fiddler。

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=flat&logo=windows)](https://www.microsoft.com/windows)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat)](LICENSE)

## ✨ 主要特性

- 🔍 **HTTP/HTTPS/WebSocket 抓包**: 支持多种协议的流量捕获与解析
- 🔐 **自动证书管理**: 自动生成并安装 CA 根证书，支持 HTTPS 中间人解密
- 🖥️ **图形用户界面**: 基于 Fyne 构建的现代化 GUI 界面
- 🎯 **系统托盘集成**: 方便的系统代理一键启用/禁用
- 🔎 **强大的搜索过滤**: 支持按 URL、Host、Method、Headers、Body 等多维度搜索
- 📊 **实时监控**: 实时显示捕获的请求和响应数据

## 📋 系统要求

- Windows 10/11
- Go 1.21 或更高版本
- 管理员权限（仅证书安装时需要）

## 🚀 快速开始

### 安装依赖

```bash
go mod download
```

### 编译运行

```bash
# 开发模式运行
go run main.go

# 编译可执行文件
go build -o packet-capture.exe
```

### 首次使用

1. 程序启动时会自动在 `%APPDATA%\PacketCaptureTool` 目录生成根证书
2. 如果证书未安装，程序会提示安装（可能需要管理员权限）
3. 在 GUI 界面设置代理端口（默认 8888）
4. 使用系统托盘菜单启用系统代理
5. 开始捕获网络流量

## 📖 使用说明

### GUI 界面功能

- **代理端口设置**: 自定义代理监听端口（默认 8888）
- **HTTPS 解密开关**: 启用/禁用 HTTPS 流量解密
- **重新安装证书**: 如果证书失效，可重新安装
- **清空请求**: 清除当前捕获的所有请求记录
- **搜索框**: 实时过滤和搜索请求

### 系统托盘功能

- ✅ **启用代理**: 将系统代理设置为 127.0.0.1:8888
- ❌ **禁用代理**: 清除系统代理设置
- 📊 **显示主窗口**: 打开主界面
- 🚪 **退出**: 安全退出程序

### 搜索功能

支持以下字段的搜索：
- URL 路径（如 `/api/user`）
- Host（如 `example.com`）
- 请求方法（GET、POST、PUT 等）
- 响应状态码（200、404 等）
- 请求头和响应头
- 请求体和响应体内容

## 📁 项目结构

```
packet-capture-tool/
├── main.go                          # 程序入口
├── go.mod / go.sum                  # Go 模块配置
├── internal/                        # 内部包
│   ├── app/                         # 应用程序主控制器
│   ├── cert/                        # 证书管理（生成、安装、签发）
│   ├── config/                      # 配置管理（持久化、线程安全）
│   ├── gui/                         # GUI 界面（Fyne）
│   ├── proxy/                       # 代理服务器（HTTP/HTTPS/WebSocket）
│   ├── storage/                     # 数据存储（搜索、过滤）
│   └── systray_manager/             # 系统托盘管理
├── docs/                            # 文档
│   ├── ARCHITECTURE.md              # 系统架构设计
│   ├── BUILD.md                     # 构建指南
│   ├── QUICKSTART.md                # 快速入门（5分钟上手）
│   └── USAGE.md                     # 详细使用说明
├── examples/                        # 示例和测试
│   └── simple_test.go              # 代理功能测试
├── build.bat                        # Windows 编译脚本
├── run-dev.bat                      # 开发模式运行
├── install-deps.bat                 # 依赖安装
├── verify-project.bat               # 项目验证
├── Makefile                         # Make 构建配置
├── CONTRIBUTING.md                  # 贡献指南
├── PROJECT_SUMMARY.md               # 项目总结
└── README.md                        # 项目说明（本文件）
```

## 🔒 安全说明

- 所有 MITM（中间人）行为仅限本机回环地址（127.0.0.1）
- 不会监听外网接口，确保安全性
- 根证书仅在本地使用，不会泄露
- 可随时在 GUI 中查看证书指纹

## 🛠️ 技术栈

- **编程语言**: Go 1.21+
- **GUI 框架**: Fyne v2
- **代理核心**: net/http + 自定义 MITM 引擎
- **系统托盘**: getlantern/systray
- **证书管理**: crypto/tls, crypto/x509
- **系统集成**: Windows Registry API

## ⚠️ 注意事项

1. **证书安装**: 首次运行时需要管理员权限安装根证书
2. **代理设置**: 使用完毕后记得禁用系统代理，避免影响正常网络访问
3. **内存限制**: 默认在内存中保留最新的 10000 条请求记录
4. **WebSocket**: 当前版本对 WebSocket 的支持有限

## 🐛 故障排除

### 证书安装失败
- 确保以管理员权限运行程序
- 手动运行：`certutil -addstore Root <证书路径>`

### 无法捕获 HTTPS 流量
- 检查证书是否正确安装
- 确认 HTTPS 解密开关已启用
- 验证系统代理设置是否正确

### 代理无法启动
- 检查端口是否被占用
- 尝试更换其他端口（如 8889、8890）

## 📝 开发计划

### v1.1 (短期)
- [ ] 完整的 WebSocket 流量解析
- [ ] 请求/响应详情查看窗口
- [ ] 请求重放功能
- [ ] 导出为 HAR 格式

### v1.5 (中期)
- [ ] 断点调试（修改请求/响应）
- [ ] 自动化规则引擎
- [ ] 流量统计和可视化图表
- [ ] 性能分析工具

### v2.0 (长期)
- [ ] 插件系统支持
- [ ] 分布式抓包
- [ ] 云端协作功能
- [ ] 移动端支持

## 📊 项目统计

- **代码量**: ~1,700 行 Go 代码
- **测试覆盖**: 配置、存储、搜索等核心模块
- **文档**: 7 个完整文档文件
- **模块**: 8 个核心模块，完全模块化设计

## 📄 许可证

本项目仅供学习和研究使用。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

在贡献之前，请阅读 [贡献指南](CONTRIBUTING.md)。

### 快速开始贡献

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

### 需要帮助的领域

- 🐛 Bug 修复
- ✨ 新功能实现
- 📝 文档改进
- 🧪 测试覆盖
- 🌍 国际化支持

## 🎓 学习资源

- [快速入门指南](docs/QUICKSTART.md) - 5分钟快速上手
- [系统架构文档](docs/ARCHITECTURE.md) - 深入了解系统设计
- [构建指南](docs/BUILD.md) - 编译和构建详解
- [使用手册](docs/USAGE.md) - 完整的使用说明
- [项目总结](PROJECT_SUMMARY.md) - 项目完整总结

## 🙏 致谢

本项目使用了以下优秀的开源项目：

- [Fyne](https://fyne.io/) - 跨平台 GUI 框架
- [getlantern/systray](https://github.com/getlantern/systray) - 系统托盘支持
- Go 标准库 - 强大的网络和加密支持

## 📞 支持

如有问题或建议：

1. 📖 查看 [文档](docs/)
2. 🔍 搜索已有的 [Issues](https://github.com/yourusername/packet-capture-tool/issues)
3. 💬 创建新的 Issue
4. 📧 联系维护者

## ⭐ Star History

如果这个项目对你有帮助，请给它一个 Star！⭐

---

**⚠️ 免责声明**: 此工具仅应用于合法的网络调试和开发测试。请勿用于任何非法用途。使用者需自行承担使用责任。

**版本**: v1.0.0 | **更新时间**: 2025-11-15 | **状态**: ✅ 稳定可用
