package route

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/littlecheny/go-backend/domain"
	"github.com/littlecheny/go-backend/logger"
)

// getBalanceHandler 查询钱包余额
// @Summary 查询钱包余额
// @Description 查询指定地址的以太坊余额（支持缓存）
// @Tags wallet
// @Security BearerAuth
// @Produce json
// @Param address path string true "钱包地址"
// @Success 200 {object} map[string]string "{address, balance}"
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /wallet/balance/{address} [get]
func getBalanceHandler(env *bootstrap.Env, app bootstrap.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")
		if address == "" {
			address = env.WalletAddress
		}

		if address == "" {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "missing wallet address"})
			return
		}

		// 构造缓存键
		cacheKey := "balance:" + address

		// 从缓存中获取余额
		if app.Redis != nil {
			balance, err := app.Redis.Get(c.Request.Context(), cacheKey).Result()
			if err == nil {
				logger.Log.WithContext(c.Request.Context()).WithField("address", address).Debug("Balance retrieved from cache")
				c.JSON(http.StatusOK, gin.H{
					"address": address,
					"balance": balance,
					"cached":  true,
				})
				return
			}
		}

		// 如果缓存中没有余额，则从以太坊节点查询
		balance, err := app.EthereumService.GetBalance(c.Request.Context(), address)
		if err != nil {
			logger.Log.WithContext(c.Request.Context()).WithField("address", address).Errorf("GetBalance error: %v", err)
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "failed to get balance"})
			return
		}

		// 缓存余额
		if app.Redis != nil {
			ttl := time.Duration(env.BalanceCacheTTLSeconds) * time.Second
			_ = app.Redis.Set(c.Request.Context(), cacheKey, balance, ttl).Err()
		}

		c.JSON(http.StatusOK, gin.H{
			"address": address,
			"balance": balance,
			"cached":  false,
		})
	}
}

func NewWalletBalanceRouter(env *bootstrap.Env, app bootstrap.Application, router *gin.RouterGroup) {
	router.GET("/balance/:address", getBalanceHandler(env, app))
	// 默认地址查询（使用配置中的地址）
	router.GET("/balance", getBalanceHandler(env, app))
}
