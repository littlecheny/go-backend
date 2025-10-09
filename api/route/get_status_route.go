package route

import (
	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/bootstrap"
)

// 聚合系统状态：应用环境、Redis 连接状态（ping）、以太坊最新区块高度
func NewGetstatusRouter(env *bootstrap.Env, app bootstrap.Application, router *gin.RouterGroup) {
	router.GET("/status", func(c *gin.Context) {
		ctx := c.Request.Context()

		redisPing := "unknown"
		if app.Redis != nil {
			if _, err := app.Redis.Ping(ctx).Result(); err == nil {
				redisPing = "ok"
			} else {
				redisPing = "error"
			}
		}

		var latestBlock uint64
		var chainErr string
		if app.EthereumService != nil {
			lb, err := app.EthereumService.GetLatestBlockNumber(ctx)
			if err == nil {
				latestBlock = lb
			} else {
				chainErr = err.Error()
			}
		}

		c.JSON(200, gin.H{
			"app": env.AppEnv,
			"redis": gin.H{
				"host": env.RedisHost,
				"ping": redisPing,
			},
			"ethereum": gin.H{
				"network":      env.DefaultNetwork,
				"latest_block": latestBlock,
				"error":        chainErr,
			},
		})
	})
}
