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
// @Security BearerAuth
// @Produce json
// @Param taskID path string true "任务ID"
// @Success 200 {object} map[string]string "{task_id, tx_hash, status}"
// @Failure 400 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
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
		
		// 查询任务状态
		status, err := app.Redis.Get(ctx, "task:status:"+taskID).Result()
		if err != nil {
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "task not found"})
			return
		}

		// 查询交易哈希（可能不存在，如果交易还未发送）
		txHash, _ := app.Redis.Get(ctx, "task:hash:"+taskID).Result()

		response := gin.H{
			"task_id": taskID,
			"status":  status,
		}

		if txHash != "" {
			response["tx_hash"] = txHash
		}

		c.JSON(http.StatusOK, response)
	}
}

// 查询任务：GET /tx/:taskID
// 返回：{ task_id, tx_hash, status }
func NewTxTaskRouter(env *bootstrap.Env, app bootstrap.Application, router *gin.RouterGroup) {
	router.GET(":taskID", getTxTaskHandler(app))
}