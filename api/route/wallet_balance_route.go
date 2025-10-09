package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/bootstrap"
)

func NewWalletBalanceRouter(env *bootstrap.Env, app bootstrap.Application, router *gin.RouterGroup) {
	router.GET("/balance", func(c *gin.Context) {
		address := env.WalletAddress
		// 从缓存中获取余额
		balance, err := app.Redis.Get(c.Request.Context(), address).Result()
		if err == nil {
			c.JSON(200, gin.H{"balance": balance})
			return
		}
		// 如果缓存中没有余额，则从以太坊节点查询
		balance, err = app.EthereumService.GetBalance(c.Request.Context(), address)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"balance": balance})
		// 缓存余额
		app.Redis.Set(c.Request.Context(), address, balance, time.Duration(env.BalanceCacheTTLSeconds)*time.Second)
	})

}
