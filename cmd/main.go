package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	route "github.com/littlecheny/go-backend/api/route"
	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/littlecheny/go-backend/logger"
)

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

	timeout := time.Duration(env.ContextTimeout) * time.Second

	r := gin.Default()

	route.Setup(env, db, r, timeout, app)

	r.Run(env.ServerAddress)
}
