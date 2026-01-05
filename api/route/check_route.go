package route

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/bootstrap"
)

// healthCheckHandler 健康检查端点
// @Summary 健康检查
// @Description 检查服务健康状态
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func healthCheckHandler(env *bootstrap.Env, app bootstrap.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "1.0.0",
			"services": gin.H{
				"database":  "connected",
				"redis":     "connected",
				"ethereum":  "connected",
			},
		}

		// 检查 MongoDB 连接
		if app.Mongo == nil {
			status["services"].(gin.H)["database"] = "disconnected"
			status["status"] = "degraded"
		}

		// 检查 Redis 连接
		if app.Redis == nil {
			status["services"].(gin.H)["redis"] = "disconnected"
			status["status"] = "degraded"
		}

		// 检查 Ethereum 服务
		if app.EthereumService == nil {
			status["services"].(gin.H)["ethereum"] = "disconnected"
			status["status"] = "degraded"
		}

		c.JSON(http.StatusOK, status)
	}
}

func NewCheckRouter(env *bootstrap.Env, app bootstrap.Application, router *gin.RouterGroup) {
	router.GET("/health", healthCheckHandler(env, app))
}
