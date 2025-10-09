package route

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/littlecheny/go-backend/domain"
	"github.com/littlecheny/go-backend/logger"
)

// getTxStatusHandler 提供 GET /tx/:hash/status 的处理函数
// @Summary 查询交易状态
// @Description 根据交易哈希查询区块链交易状态
// @Tags tx
// @Produce json
// @Param hash path string true "交易哈希"
// @Success 200 {object} map[string]string "{hash, status}"
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /tx/{hash}/status [get]
func getTxStatusHandler(app bootstrap.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		if hash == "" {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "missing tx hash"})
			return
		}

		status, err := app.EthereumService.GetTransactionStatus(c.Request.Context(), hash)
		if err != nil {
			logger.Log.WithContext(c.Request.Context()).WithField("hash", hash).Errorf("GetTransactionStatus error: %v", err)
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"hash":   hash,
			"status": status,
		})
	}
}

// 交易状态查询：GET /tx/:hash/status
func NewTxStatusRouter(env *bootstrap.Env, app bootstrap.Application, router *gin.RouterGroup) {
	router.GET("/:hash/status", getTxStatusHandler(app))
}