package main

// @title Go Backend API
// @version 1.0
// @description Go backend with Ethereum tx queue and REST endpoints.
// @schemes http https
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	route "github.com/littlecheny/go-backend/api/route"
	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/littlecheny/go-backend/domain"
	"github.com/littlecheny/go-backend/logger"
	docs "github.com/littlecheny/go-backend/docs"
)

func startTxWorker(ctx context.Context, app bootstrap.Application) {
	if app.Redis == nil || app.EthereumService == nil {
		logger.Log.WithContext(ctx).Warn("TxWorker not started: Redis or EthereumService unavailable")
		return
	}
	go func() {
		for {
			// 阻塞弹出队列任务（FIFO：LPUSH + BRPOP）
			res, err := app.Redis.BRPop(ctx, 0, "tx:queue").Result()
			if err != nil {
				logger.Log.WithContext(ctx).Errorf("BRPop error: %v", err)
				time.Sleep(time.Second)
				continue
			}
			if len(res) < 2 {
				continue
			}
			data := res[1]

			var payload struct {
				TaskID string                        `json:"task_id"`
				Req    domain.TransactionSendRequest `json:"req"`
			}
			if err := json.Unmarshal([]byte(data), &payload); err != nil {
				logger.Log.WithContext(ctx).Errorf("unmarshal payload error: %v", err)
				continue
			}

			// 更新任务状态：processing
			_ = app.Redis.Set(ctx, "task:status:"+payload.TaskID, "processing", 24*time.Hour).Err()

			// 发送交易
			txHash, err := app.EthereumService.SendTransaction(ctx, &payload.Req)
			if err != nil {
				_ = app.Redis.Set(ctx, "task:status:"+payload.TaskID, "failed_to_send", 24*time.Hour).Err()
				logger.Log.WithContext(ctx).Errorf("SendTransaction error: %v", err)
				continue
			}

			// 写入交易哈希与状态
			_ = app.Redis.Set(ctx, "task:hash:"+payload.TaskID, txHash, 24*time.Hour).Err()
			_ = app.Redis.Set(ctx, "task:status:"+payload.TaskID, "sent", 24*time.Hour).Err()

			// 异步轮询回执，带退避
			go func(taskID, hash string) {
				backoff := time.Second
				maxBackoff := 30 * time.Second
				attempts := 0
				for {
					attempts++
					status, err := app.EthereumService.GetTransactionStatus(ctx, hash)
					if err == nil && status != "" {
						_ = app.Redis.Set(ctx, "task:status:"+taskID, status, 24*time.Hour).Err()
						break
					}
					if attempts >= 60 { // 最多轮询约 ~10分钟（指数退避）
						_ = app.Redis.Set(ctx, "task:status:"+taskID, "pending", 24*time.Hour).Err()
						break
					}
					time.Sleep(backoff)
					if backoff < maxBackoff {
						backoff *= 2
						if backoff > maxBackoff {
							backoff = maxBackoff
						}
					}
				}
			}(payload.TaskID, txHash)
		}
	}()
}

func main() {
	app := bootstrap.App()

	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	//初始化日志
	logger.Init()

	//初始化区块链服务
	err := bootstrap.InitBlockchain(context.Background(), env, &app)
	if err != nil {
		logger.Log.Fatalf("Failed to initialize Blockchain service: %v", err)
	}

	//初始化Redis
	app.Redis = bootstrap.NewRedisClient(env)
	defer bootstrap.CloseRedisConnection(app.Redis)

	// 启动后台交易队列Worker
	startTxWorker(context.Background(), app)

	timeout := time.Duration(env.ContextTimeout) * time.Second

	r := gin.Default()

	// Swagger 元数据配置（需在路由注册前设置）
	docs.SwaggerInfo.Title = "Go Backend API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "Go backend with Ethereum tx queue and REST endpoints."
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	route.Setup(env, db, r, timeout, app)

	r.Run(env.ServerAddress)
}
