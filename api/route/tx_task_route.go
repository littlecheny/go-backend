package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/littlecheny/go-backend/domain"
)

// getTxTaskHandler 提供 GET /tx/:taskID 的处理函数
// @Summary 查询交易任务状态
// @Description 根据 taskID 查询交易任务的当前状态与交易哈希
// @Tags tx
// @Produce json
// @Param taskID path string true "任务ID"
// @Success 200 {object} map[string]string "{task_id, tx_hash, status}"
// @Failure 400 {object} domain.ErrorResponse
// @Failure 503 {object} domain.ErrorResponse
// @Router /tx/{taskID} [get]
func getTxTaskHandler(app bootstrap.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("taskID")
		if taskID == "" {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "missing taskID"})
			return
		}

		if app.Redis == nil {
			c.JSON(http.StatusServiceUnavailable, domain.ErrorResponse{Message: "redis unavailable"})
			return
		}

		ctx := c.Request.Context()
		txHash, _ := app.Redis.Get(ctx, "task:hash:"+taskID).Result()
		status, _ := app.Redis.Get(ctx, "task:status:"+taskID).Result()

		c.JSON(http.StatusOK, gin.H{
			"task_id": taskID,
			"tx_hash": txHash,
			"status":  status,
		})
	}
}

// 查询任务：GET /tx/:taskID
// 返回：{ task_id, tx_hash, status }
func NewTxTaskRouter(env *bootstrap.Env, app bootstrap.Application, router *gin.RouterGroup) {
	router.GET(":taskID", getTxTaskHandler(app))
}