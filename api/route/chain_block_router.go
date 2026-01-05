package route

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/littlecheny/go-backend/domain"
	"github.com/littlecheny/go-backend/logger"
)

// getLatestBlockHandler 查询最新区块高度
// @Summary 查询最新区块高度
// @Description 查询区块链的最新区块高度（支持缓存）
// @Tags chain
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "{block_number, cached}"
// @Failure 500 {object} domain.ErrorResponse
// @Router /chain/block/latest [get]
func getLatestBlockHandler(env *bootstrap.Env, app bootstrap.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		cacheKey := "latest_block"

		// 从缓存中获取最新区块高度
		if app.Redis != nil {
			cachedBlock, err := app.Redis.Get(ctx, cacheKey).Result()
			if err == nil {
				blockNumber, _ := strconv.ParseUint(cachedBlock, 10, 64)
				logger.Log.WithContext(ctx).Debug("Block number retrieved from cache")
				c.JSON(http.StatusOK, gin.H{
					"block_number": blockNumber,
					"cached":       true,
					"network":      env.DefaultNetwork,
				})
				return
			}
		}

		// 从区块链查询最新区块
		latestBlock, err := app.EthereumService.GetLatestBlockNumber(ctx)
		if err != nil {
			logger.Log.WithContext(ctx).Errorf("GetLatestBlockNumber error: %v", err)
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
				Message: "failed to get latest block number",
			})
			return
		}

		// 缓存最新区块高度
		if app.Redis != nil {
			ttl := time.Duration(env.BlockCacheTTLSeconds) * time.Second
			_ = app.Redis.Set(ctx, cacheKey, latestBlock, ttl).Err()
		}

		c.JSON(http.StatusOK, gin.H{
			"block_number": latestBlock,
			"cached":       false,
			"network":      env.DefaultNetwork,
		})
	}
}

func NewChainBlockRouter(env *bootstrap.Env, app bootstrap.Application, router *gin.RouterGroup) {
	router.GET("/block/latest", getLatestBlockHandler(env, app))
}
