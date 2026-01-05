.PHONY: help build run test clean docker-build docker-up docker-down swagger deps lint

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
help: ## 显示帮助信息
	@echo "可用命令："
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# 安装依赖
deps: ## 安装项目依赖
	@echo "正在安装依赖..."
	@go mod download
	@go mod tidy

# 生成 Swagger 文档
swagger: ## 生成 Swagger API 文档
	@echo "正在生成 Swagger 文档..."
	@swag init -g cmd/main.go
	@echo "Swagger 文档生成完成"

# 代码检查
lint: ## 运行代码检查
	@echo "正在运行代码检查..."
	@go fmt ./...
	@go vet ./...

# 构建
build: ## 构建应用程序
	@echo "正在构建应用..."
	@go build -o bin/server cmd/main.go
	@echo "构建完成: bin/server"

# 运行
run: ## 运行应用程序
	@echo "正在启动应用..."
	@go run cmd/main.go

# 测试
test: ## 运行单元测试
	@echo "正在运行测试..."
	@go test -v -cover ./...

# 测试覆盖率
test-coverage: ## 生成测试覆盖率报告
	@echo "正在生成测试覆盖率报告..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 生成 Mock
mock: ## 生成 Mock 文件
	@echo "正在生成 Mock 文件..."
	@mockery --all --dir=domain --output=domain/mocks
	@echo "Mock 文件生成完成"

# 清理
clean: ## 清理构建产物
	@echo "正在清理..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "清理完成"

# Docker 构建
docker-build: ## 构建 Docker 镜像
	@echo "正在构建 Docker 镜像..."
	@docker build -t go-backend:latest .
	@echo "Docker 镜像构建完成"

# Docker 启动
docker-up: ## 启动 Docker 容器
	@echo "正在启动 Docker 容器..."
	@docker-compose up -d
	@echo "Docker 容器已启动"

# Docker 停止
docker-down: ## 停止 Docker 容器
	@echo "正在停止 Docker 容器..."
	@docker-compose down
	@echo "Docker 容器已停止"

# Docker 查看日志
docker-logs: ## 查看 Docker 容器日志
	@docker-compose logs -f api

# 开发环境（启动依赖服务）
dev-up: ## 启动开发环境依赖服务（MongoDB + Redis）
	@echo "正在启动开发环境依赖服务..."
	@docker-compose up -d mongo redis
	@echo "依赖服务已启动"

# 完整部署
deploy: deps swagger build docker-build docker-up ## 完整部署流程
	@echo "部署完成！"

# 本地开发
dev: dev-up run ## 启动本地开发环境

