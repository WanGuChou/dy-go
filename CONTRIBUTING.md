# 贡献指南

感谢你考虑为 Packet Capture Tool 做出贡献！

## 如何贡献

### 报告 Bug

如果你发现了 bug，请创建一个 Issue，包含以下信息：

1. **Bug 描述**: 简短描述问题
2. **重现步骤**: 详细的重现步骤
3. **预期行为**: 你期望发生什么
4. **实际行为**: 实际发生了什么
5. **环境信息**:
   - 操作系统版本（Windows 10/11）
   - Go 版本
   - 程序版本
6. **日志输出**: 相关的日志或错误信息
7. **截图**: 如果适用，提供截图

### 提出功能请求

如果你有新功能的想法：

1. 检查是否已有类似的 Issue
2. 创建新 Issue，使用 "Feature Request" 标签
3. 详细描述功能需求和使用场景
4. 如果可能，提供实现思路

### 提交代码

1. **Fork 项目**
   ```bash
   # 在 GitHub 上点击 Fork 按钮
   git clone https://github.com/your-username/packet-capture-tool.git
   cd packet-capture-tool
   ```

2. **创建功能分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或
   git checkout -b fix/bug-description
   ```

3. **编写代码**
   - 遵循项目的代码风格
   - 添加必要的注释
   - 编写单元测试

4. **运行测试**
   ```bash
   go test ./...
   go vet ./...
   go fmt ./...
   ```

5. **提交更改**
   ```bash
   git add .
   git commit -m "Add: your feature description"
   git push origin feature/your-feature-name
   ```

6. **创建 Pull Request**
   - 在 GitHub 上创建 PR
   - 详细描述你的更改
   - 链接相关的 Issue

## 代码规范

### Go 代码风格

遵循 [Effective Go](https://golang.org/doc/effective_go) 和 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)。

**基本规则**:

```go
// ✅ 好的示例
func ProcessRequest(req *Request) error {
    if req == nil {
        return errors.New("request is nil")
    }
    
    // 处理逻辑
    return nil
}

// ❌ 不好的示例
func process_request(req *Request) error {  // 应使用驼峰命名
    if req!=nil{  // 应有空格
        //处理逻辑  // 应有空格
        return nil
    }
    return errors.New("request is nil")
}
```

### 命名规范

- **包名**: 小写，单个单词（如 `proxy`, `storage`）
- **文件名**: 小写，下划线分隔（如 `cert_manager.go`）
- **类型名**: 大驼峰（如 `CertManager`）
- **函数名**: 大驼峰（公开）或小驼峰（私有）
- **变量名**: 小驼峰
- **常量名**: 大驼峰或全大写下划线

### 注释规范

```go
// Package proxy 提供 HTTP/HTTPS 代理服务器功能
package proxy

// Server 代理服务器结构体
// 负责监听端口、拦截请求、转发流量
type Server struct {
    port int
    // ... 其他字段
}

// NewServer 创建新的代理服务器实例
// 参数:
//   - port: 监听端口
//   - certManager: 证书管理器
// 返回:
//   - *Server: 服务器实例
func NewServer(port int, certManager *cert.Manager) *Server {
    // ...
}
```

### 错误处理

```go
// ✅ 好的示例
result, err := someOperation()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// ❌ 不好的示例
result, _ := someOperation()  // 不要忽略错误
```

### 测试规范

```go
// 测试文件命名: xxx_test.go
package storage

import "testing"

// 测试函数命名: TestXxx
func TestStorage_Add(t *testing.T) {
    // Arrange
    s := NewStorage()
    req := &CapturedRequest{...}
    
    // Act
    s.Add(req)
    
    // Assert
    if s.Count() != 1 {
        t.Errorf("Expected count 1, got %d", s.Count())
    }
}
```

## 项目结构

```
packet-capture-tool/
├── main.go              # 程序入口
├── internal/            # 内部包
│   ├── app/            # 应用控制器
│   ├── cert/           # 证书管理
│   ├── config/         # 配置管理
│   ├── gui/            # GUI 界面
│   ├── proxy/          # 代理服务器
│   ├── storage/        # 数据存储
│   └── systray_manager/ # 系统托盘
├── docs/               # 文档
├── examples/           # 示例代码
└── README.md           # 项目说明
```

## 开发流程

### 1. 设置开发环境

```bash
# 安装依赖
go mod download

# 安装开发工具
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 2. 开发新功能

1. 在 `internal/` 下创建或修改相应模块
2. 编写单元测试
3. 更新文档
4. 运行所有测试和检查

### 3. 代码检查

```bash
# 格式化代码
go fmt ./...

# 静态检查
go vet ./...

# Lint 检查（如果已安装 golangci-lint）
golangci-lint run

# 运行测试
go test -v ./...
go test -race ./...  # 检测并发问题
```

### 4. 提交前检查清单

- [ ] 代码已格式化（`go fmt`）
- [ ] 通过静态检查（`go vet`）
- [ ] 所有测试通过（`go test`）
- [ ] 添加了必要的注释
- [ ] 更新了相关文档
- [ ] 提交信息清晰明确

## 提交信息规范

使用语义化的提交信息：

```
类型: 简短描述

详细描述（可选）

关闭的 Issue（可选）
```

**类型**:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具相关

**示例**:

```
feat: 添加请求重放功能

实现了捕获请求的重放功能，用户可以重新发送之前捕获的请求。

Closes #123
```

## Pull Request 流程

1. **创建 PR**: 从你的 fork 创建 PR 到主仓库
2. **描述更改**: 详细说明你做了什么以及为什么
3. **链接 Issue**: 如果解决了某个 Issue，添加 "Closes #issue_number"
4. **等待审查**: 维护者会审查你的代码
5. **修改反馈**: 根据审查意见进行修改
6. **合并**: 审查通过后，维护者会合并你的 PR

## 代码审查标准

审查者会检查：

1. **功能性**: 代码是否正确实现了预期功能
2. **代码质量**: 是否遵循最佳实践和项目规范
3. **测试**: 是否有充分的测试覆盖
4. **文档**: 是否更新了相关文档
5. **性能**: 是否有明显的性能问题
6. **安全性**: 是否引入了安全风险

## 需要帮助的领域

我们特别欢迎以下方面的贡献：

- 🐛 **Bug 修复**: 解决已知问题
- ✨ **新功能**: 实现计划中的功能
- 📝 **文档**: 改进文档和示例
- 🧪 **测试**: 增加测试覆盖率
- 🌍 **国际化**: 添加多语言支持
- 🎨 **UI/UX**: 改进用户界面
- ⚡ **性能**: 优化性能

## 行为准则

### 我们的承诺

为了营造一个开放和友好的环境，我们承诺：

- 尊重不同的观点和经验
- 接受建设性的批评
- 关注对社区最有利的事情
- 对其他社区成员保持同理心

### 不可接受的行为

- 使用性暗示的语言或图像
- 人身攻击、侮辱或贬低性评论
- 公开或私下的骚扰
- 未经许可发布他人的隐私信息
- 其他在专业环境中不适当的行为

## 许可证

通过贡献代码，你同意你的贡献将按照项目的许可证进行许可。

## 联系方式

如有疑问，请：

1. 查看文档
2. 搜索已有的 Issue
3. 创建新的 Issue
4. 发送邮件到维护者

## 感谢

感谢所有为这个项目做出贡献的人！

---

再次感谢你的贡献！ 🎉
