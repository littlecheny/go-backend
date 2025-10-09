package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/bootstrap"
)

func NewChainBlockRouter(env *bootstrap.Env, app bootstrap.Application, router *gin.RouterGroup) {
	router.GET("/latest-block", func(c *gin.Context) {
		// 从缓存中获取最新区块高度
		latestBlock, err := app.Redis.Get(c.Request.Context(), "latest_block").Uint64()
		if err == nil {
			c.JSON(200, gin.H{"latest_block": latestBlock})
			return
		}

		latestBlock, err = app.EthereumService.GetLatestBlockNumber(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{
				"message": "internal server error",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "ok",
			"data":    latestBlock,
		})
		// 缓存最新区块高度
		app.Redis.Set(c.Request.Context(), "latest_block", latestBlock, time.Duration(env.BlockCacheTTLSeconds)*time.Second)
	})
}
