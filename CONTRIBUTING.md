# 贡献指南

感谢你考虑为本项目做出贡献！

## 开发流程

### 1. Fork 项目

点击 GitHub 页面右上角的 "Fork" 按钮，创建一个你自己的副本。

### 2. 克隆仓库

```bash
git clone https://github.com/YOUR_USERNAME/go-backend.git
cd go-backend
```

### 3. 创建分支

```bash
git checkout -b feature/your-feature-name
```

分支命名规范：
- `feature/xxx` - 新功能
- `bugfix/xxx` - Bug 修复
- `hotfix/xxx` - 紧急修复
- `docs/xxx` - 文档更新
- `refactor/xxx` - 重构
- `test/xxx` - 测试相关

### 4. 进行更改

遵循项目的代码规范：
- 使用 `go fmt` 格式化代码
- 使用 `go vet` 检查代码
- 添加必要的单元测试
- 更新相关文档

### 5. 提交更改

提交信息格式：

```
<type>: <subject>

<body>

<footer>
```

Type 类型：
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具链相关

示例：

```bash
git commit -m "feat: 添加用户头像上传功能

实现了用户头像上传和存储功能：
- 支持 jpg、png 格式
- 自动压缩和裁剪
- 使用七牛云存储

Closes #123"
```

### 6. 推送到 GitHub

```bash
git push origin feature/your-feature-name
```

### 7. 创建 Pull Request

1. 前往你 Fork 的仓库页面
2. 点击 "New Pull Request" 按钮
3. 填写 PR 标题和描述
4. 等待代码审查

## 代码规范

### Go 代码风格

- 遵循 [Effective Go](https://golang.org/doc/effective_go.html) 指南
- 使用 `gofmt` 格式化代码
- 使用有意义的变量和函数名
- 添加必要的注释，特别是导出的函数和类型
- 避免使用全局变量

### 命名规范

**变量命名**
```go
// 好的命名
userEmail := "user@example.com"
maxRetryCount := 3

// 不好的命名
e := "user@example.com"
n := 3
```

**函数命名**
```go
// 好的命名
func GetUserByEmail(email string) (User, error)
func ValidatePassword(password string) bool

// 不好的命名
func get(e string) (User, error)
func check(p string) bool
```

### 错误处理

```go
// 好的做法
result, err := SomeFunction()
if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}

// 不好的做法
result, _ := SomeFunction()  // 忽略错误
```

### 测试

- 为新功能编写单元测试
- 测试函数命名：`TestFunctionName_Scenario`
- 使用 table-driven tests

```go
func TestGetUserByEmail_Success(t *testing.T) {
    // 测试代码
}

func TestGetUserByEmail_NotFound(t *testing.T) {
    // 测试代码
}
```

## 文档规范

### 代码注释

```go
// GetUserByEmail 根据邮箱地址查询用户
// 如果用户不存在，返回 ErrUserNotFound 错误
func GetUserByEmail(email string) (User, error) {
    // 实现代码
}
```

### Swagger 注释

```go
// @Summary 用户登录
// @Description 使用邮箱和密码登录，返回 JWT Token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Router /login [post]
func (lc *LoginController) Login(c *gin.Context) {
    // 实现代码
}
```

## 测试指南

### 运行测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./repository -v

# 运行特定测试
go test -run TestGetUserByEmail ./repository -v

# 生成覆盖率报告
make test-coverage
```

### 编写测试

```go
func TestUserRepository_Create(t *testing.T) {
    // Arrange
    mockDB := new(mocks.Database)
    repo := NewUserRepository(mockDB, "users")
    user := &domain.User{
        Name:  "Test User",
        Email: "test@example.com",
    }

    // Act
    err := repo.Create(context.Background(), user)

    // Assert
    assert.NoError(t, err)
    mockDB.AssertExpectations(t)
}
```

## Pull Request 检查清单

提交 PR 前，请确保：

- [ ] 代码已经过 `go fmt` 格式化
- [ ] 代码已通过 `go vet` 检查
- [ ] 已添加必要的单元测试
- [ ] 所有测试都通过
- [ ] 已更新相关文档
- [ ] 提交信息清晰明确
- [ ] PR 描述详细说明了更改内容

## 代码审查

我们会尽快审查你的 PR，可能会提出一些修改建议。请保持耐心，并积极回应审查意见。

## 问题反馈

如果你发现了 Bug 或有功能建议，请：

1. 在 GitHub Issues 中搜索，确保问题未被报告
2. 创建新的 Issue，使用合适的模板
3. 提供详细的复现步骤或用例说明

## 联系方式

如有任何问题，欢迎通过以下方式联系：

- GitHub Issues: [项目 Issues](https://github.com/littlecheny/go-backend/issues)
- Email: your-email@example.com

## 行为准则

- 尊重他人
- 接受建设性批评
- 关注对项目最有利的事情
- 对社区其他成员表示同理心

感谢你的贡献！🎉

