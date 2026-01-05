# Go Backend - 企业级区块链交易系统

![Go Version](https://img.shields.io/badge/Go-1.24.2-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)

一个基于 Go 语言开发的企业级后端系统，集成 Ethereum 区块链交易处理、异步任务队列、JWT 认证等功能的 RESTful API 服务。

## 📋 目录

- [项目架构](#项目架构)
- [技术亮点](#技术亮点)
- [核心功能](#核心功能)
- [技术栈](#技术栈)
- [系统架构](#系统架构)
- [快速开始](#快速开始)
- [API 文档](#api-文档)
- [项目结构](#项目结构)
- [性能优化](#性能优化)
- [安全性](#安全性)

## 🏗️ 项目架构

本项目采用 **Clean Architecture（清晰架构）** 设计模式，实现了高内聚低耦合的代码结构：

```
┌─────────────────────────────────────────────────────────┐
│                     API Layer (Gin)                      │
│              Controllers + Middleware + Routes            │
└───────────────────────┬─────────────────────────────────┘
                        │
┌───────────────────────▼─────────────────────────────────┐
│                   Use Case Layer                         │
│              Business Logic + Orchestration              │
└───────────────────────┬─────────────────────────────────┘
                        │
┌───────────────────────▼─────────────────────────────────┐
│                   Domain Layer                           │
│          Entities + Interfaces + Domain Logic           │
└───────────────────────┬─────────────────────────────────┘
                        │
┌───────────────────────▼─────────────────────────────────┐
│               Repository + Service Layer                 │
│         MongoDB Repository + Ethereum Service            │
└─────────────────────────────────────────────────────────┘
```

### 架构优势

- **依赖倒置原则（DIP）**：核心业务逻辑不依赖外部实现
- **接口隔离**：通过 interface 定义清晰的契约
- **单一职责**：每层职责明确，易于维护和测试
- **可测试性**：完整的 Mock 支持，便于单元测试

## 🌟 技术亮点

### 1. 异步任务队列架构

采用 **Redis List + Worker Pool** 模式处理区块链交易：

- **FIFO 队列**：使用 Redis LPUSH/BRPOP 实现可靠的任务队列
- **幂等性保证**：通过 `Idempotency-Key` 机制防止重复提交
- **异步处理**：立即返回任务 ID（202 Accepted），后台异步处理
- **智能轮询**：指数退避算法查询交易状态，最多轮询 10 分钟
- **状态追踪**：完整的任务生命周期管理（queued → processing → sent → success/failed）

```go
// 任务状态流转
queued → processing → sent → [success/failed/pending]
```

### 2. 完善的中间件体系

- **JWT 认证中间件**：基于 Bearer Token 的无状态认证
- **限流中间件**：基于内存的高性能限流器（每分钟 100 次请求）
- **CORS 中间件**：跨域资源共享支持
- **错误处理中间件**：统一的错误响应格式
- **日志中间件**：结构化日志记录（Logrus）

### 3. 区块链交易处理

- **多链支持**：当前支持 Ethereum Sepolia 测试网，架构可扩展至其他 EVM 链
- **私钥安全**：使用环境变量管理私钥，支持 Hardware Wallet 集成
- **Gas 优化**：可配置 Gas Price 和 Gas Limit
- **交易签名**：使用 go-ethereum 标准库进行 EIP-155 签名
- **状态查询**：实时查询交易回执和状态

### 4. 高性能缓存策略

- **Redis 缓存**：余额和区块高度信息缓存
- **TTL 管理**：可配置的缓存过期时间
- **缓存预热**：应用启动时预加载关键数据

### 5. 完整的 API 文档

- **Swagger/OpenAPI 3.0**：自动生成交互式 API 文档
- **注释驱动**：通过代码注释自动生成文档
- **在线测试**：支持在 Swagger UI 中直接测试 API

## 🚀 核心功能

### 用户认证系统
- 用户注册/登录
- JWT Access Token + Refresh Token 双令牌机制
- 密码加密存储（Bcrypt）
- 邮箱格式验证

### 区块链交易管理
- 异步发送以太坊交易
- 交易状态实时查询
- 任务队列管理
- 交易历史追踪

### 账户管理
- 查询钱包余额（支持缓存）
- 查询最新区块高度
- 多地址支持

### 系统监控
- 健康检查端点
- 系统状态监控
- 日志追踪

## 🛠️ 技术栈

### 核心框架
- **Gin**: 高性能 HTTP Web 框架
- **Go Ethereum**: 官方以太坊 Go 客户端库

### 数据存储
- **MongoDB**: 文档型数据库，存储用户信息
- **Redis**: 内存数据库，用于任务队列和缓存

### 安全认证
- **JWT**: 无状态认证方案
- **Bcrypt**: 密码哈希算法

### 开发工具
- **Swagger**: API 文档自动生成
- **Viper**: 配置管理
- **Logrus**: 结构化日志
- **Testify**: 单元测试框架
- **Mockery**: Mock 生成工具

### DevOps
- **Docker**: 容器化部署
- **Docker Compose**: 多容器编排

## 📊 系统架构

### 数据流图

```
Client Request
      ↓
  Gin Router (Middleware Pipeline)
      ↓
  Controller (Request Validation)
      ↓
  Use Case (Business Logic)
      ↓
  Repository/Service (Data Access)
      ↓
  MongoDB / Redis / Ethereum
```

### 异步交易处理流程

```
POST /tx/send (Client)
      ↓
  Validate Request + Idempotency Check
      ↓
  Generate Task ID → Return 202 Accepted
      ↓
  LPUSH to Redis Queue (tx:queue)
      ↓
Background Worker (BRPOP)
      ↓
  Sign & Send Transaction to Ethereum
      ↓
  Update Status: sent
      ↓
  Async Poll Receipt (Exponential Backoff)
      ↓
  Update Status: success/failed
```

## 🏃 快速开始

### 前置要求

- Go 1.24.2+
- Docker & Docker Compose
- MongoDB 6.0+
- Redis 7+

### 本地开发

1. **克隆项目**
```bash
git clone https://github.com/littlecheny/go-backend.git
cd go-backend
```

2. **配置环境变量**
```bash
# 从模板创建 .env 文件
cp .env.example .env

# 编辑 .env 文件，填入实际配置
vim .env
```

3. **启动依赖服务**
```bash
docker-compose up -d mongo redis
```

4. **安装依赖**
```bash
go mod download
```

5. **生成 Swagger 文档**
```bash
# 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g cmd/main.go
```

6. **运行服务**
```bash
go run cmd/main.go
```

服务将在 `http://localhost:8080` 启动

### Docker 部署

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f api

# 停止服务
docker-compose down
```

## 📖 API 文档

启动服务后，访问 Swagger UI：

```
http://localhost:8080/swagger/index.html
```

### 主要端点

#### 认证相关
- `POST /signup` - 用户注册
- `POST /login` - 用户登录

#### 区块链交易（需要认证）
- `POST /tx/send` - 发送交易（异步）
- `GET /tx/status/:taskId` - 查询任务状态
- `GET /tx/:txHash` - 查询交易详情

#### 账户管理（需要认证）
- `GET /wallet/balance/:address` - 查询钱包余额
- `GET /chain/block/latest` - 查询最新区块

#### 系统监控
- `GET /health` - 健康检查
- `GET /status` - 系统状态

### 请求示例

**发送交易**
```bash
curl -X POST http://localhost:8080/tx/send \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Idempotency-Key: unique-request-id-123" \
  -H "Content-Type: application/json" \
  -d '{
    "from": "0xYourAddress",
    "to": "0xRecipientAddress",
    "value": "1000000000000000000",
    "gas_price": "20000000000",
    "gas_limit": "21000"
  }'
```

**查询任务状态**
```bash
curl http://localhost:8080/tx/status/20240105T120000.000Z00:00 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📁 项目结构

```
go-backend/
├── api/
│   ├── controller/          # 控制器层
│   ├── middleware/          # 中间件
│   └── route/              # 路由配置
├── bootstrap/              # 应用初始化
│   ├── app.go             # 应用实例
│   ├── database.go        # MongoDB 连接
│   ├── redis.go           # Redis 连接
│   └── blockchain.go      # 区块链服务初始化
├── cmd/
│   └── main.go            # 应用入口
├── domain/                # 领域层（核心业务逻辑）
│   ├── user.go            # 用户实体
│   ├── task.go            # 任务实体
│   ├── blockchain.go      # 区块链接口定义
│   └── mocks/             # Mock 实现
├── repository/            # 数据访问层
│   ├── user_repository.go
│   └── task_repository.go
├── usecase/               # 用例层（业务编排）
│   ├── login_usecase.go
│   └── signup_usecase.go
├── services/              # 外部服务
│   ├── ethereum_sepolia_service.go
│   └── ethereum_mock_service.go
├── internal/              # 内部工具包
│   └── tokenutil/         # JWT 工具
├── logger/                # 日志配置
├── mongo/                 # MongoDB 抽象层
├── docs/                  # Swagger 文档
├── docker-compose.yml     # Docker 编排
├── Dockerfile             # 容器镜像
├── go.mod                 # Go 模块定义
└── README.md              # 项目文档
```

### 目录职责说明

- **api/**: API 层，处理 HTTP 请求
- **domain/**: 领域层，定义核心业务实体和接口
- **usecase/**: 用例层，编排业务流程
- **repository/**: 仓储层，数据持久化
- **services/**: 服务层，外部系统集成
- **bootstrap/**: 引导层，应用启动和依赖注入

## ⚡ 性能优化

### 1. 数据库优化
- MongoDB 索引优化（邮箱字段唯一索引）
- 连接池复用
- 查询字段投影（避免返回敏感字段）

### 2. 缓存策略
- Redis 缓存热点数据
- 可配置的 TTL
- 缓存穿透保护

### 3. 并发处理
- Goroutine 异步处理任务
- Worker Pool 模式
- 通道（Channel）通信

### 4. API 性能
- 请求参数验证
- 响应数据压缩
- HTTP/2 支持

## 🔒 安全性

### 1. 认证授权
- JWT 双令牌机制（Access Token + Refresh Token）
- Token 过期时间控制
- Bearer Token 认证方案

### 2. 数据安全
- 密码 Bcrypt 加密存储（Cost 10）
- 私钥环境变量管理
- HTTPS 传输加密（生产环境）

### 3. API 安全
- CORS 配置
- 限流防护
- 请求参数验证
- SQL/NoSQL 注入防护

### 4. 幂等性保证
- Idempotency-Key 机制
- 24 小时幂等窗口
- Redis 存储幂等映射

## 🧪 测试

### 运行单元测试
```bash
go test ./... -v -cover
```

### 生成测试覆盖率报告
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Mock 生成
```bash
# 使用 mockery 生成 Mock
mockery --all --dir=domain --output=domain/mocks
```

## 📈 监控与日志

### 日志级别
- **Info**: 正常操作日志
- **Warn**: 警告信息（可恢复错误）
- **Error**: 错误信息（需要关注）
- **Fatal**: 致命错误（程序退出）

### 结构化日志示例
```go
logger.Log.WithContext(ctx).WithFields(logrus.Fields{
    "task_id": taskID,
    "tx_hash": txHash,
    "status": "sent",
}).Info("Transaction sent successfully")
```

## 🚧 未来规划

- [ ] 支持多链（BSC, Polygon, Arbitrum）
- [ ] WebSocket 实时推送交易状态
- [ ] 交易批处理优化
- [ ] GraphQL API 支持
- [ ] Prometheus + Grafana 监控
- [ ] ELK 日志聚合
- [ ] Kubernetes 部署配置
- [ ] CI/CD Pipeline（GitHub Actions）
- [ ] 单元测试覆盖率提升至 80%+

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📝 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 👨‍💻 作者

**littlecheny**

- GitHub: [@littlecheny](https://github.com/littlecheny)

## 🙏 致谢

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Go Ethereum](https://github.com/ethereum/go-ethereum)
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [go-redis](https://github.com/redis/go-redis)

---

⭐ 如果这个项目对你有帮助，欢迎 Star！

