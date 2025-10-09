package route

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/littlecheny/go-backend/domain"
	"github.com/littlecheny/go-backend/logger"
)

// sendTxHandler 提供 POST /tx/send 的处理函数
// @Summary 发送以太坊交易（入队）
// @Description 通过 Redis 队列异步发送交易，支持 Idempotency-Key 幂等；立即返回 task_id
// @Tags tx
// @Accept json
// @Produce json
// @Param Idempotency-Key header string true "幂等键"
// @Param request body domain.TransactionSendRequest true "交易请求体"
// @Success 202 {object} map[string]string "{task_id, status}"
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /tx/send [post]
func sendTxHandler(app bootstrap.Application, env *bootstrap.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req domain.TransactionSendRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid request body"})
			return
		}

		idempKey := c.GetHeader("Idempotency-Key")
		if idempKey == "" {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "missing Idempotency-Key header"})
			return
		}

		// 幂等性：若已存在相同 key 的任务，则直接返回已有任务信息
		if app.Redis != nil {
			if val, err := app.Redis.Get(c.Request.Context(), "idemp:"+idempKey).Result(); err == nil && val != "" {
				c.JSON(http.StatusOK, gin.H{"task_id": val, "message": "already accepted"})
				return
			}
		}

		// 生成任务ID（可用雪花算法或 UUID，这里用时间戳简化）
		taskID := time.Now().UTC().Format("20060102T150405.000Z07:00")

		// 入队：将交易请求序列化为 JSON 放入 Redis List（队列）
		if app.Redis == nil {
			c.JSON(http.StatusServiceUnavailable, domain.ErrorResponse{Message: "queue unavailable"})
			return
		}

		payload := struct {
			TaskID string                         `json:"task_id"`
			Req    domain.TransactionSendRequest   `json:"req"`
		}{TaskID: taskID, Req: req}
		payloadBytes, _ := json.Marshal(payload)
		if err := app.Redis.LPush(c.Request.Context(), "tx:queue", payloadBytes).Err(); err != nil {
			logger.Log.WithContext(c.Request.Context()).Errorf("enqueue tx error: %v", err)
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "enqueue failed"})
			return
		}

		// 写入幂等性映射
		_ = app.Redis.Set(c.Request.Context(), "idemp:"+idempKey, taskID, 24*time.Hour).Err()

		c.JSON(http.StatusAccepted, gin.H{"task_id": taskID, "status": "queued"})
	}
}

// 交易发送：POST /tx/send
// 设计要点：
// - 使用请求头 Idempotency-Key 保证幂等。
// - 将请求写入 Redis 队列（列表）交由后台 Worker 处理签名与广播。
// - 立即返回 202 Accepted 和一个任务ID，后续可用 /tx/:taskID 查询状态。
// - 当前提交实现队列入列；后台 Worker 可后续补充。
func NewTxSendRouter(env *bootstrap.Env, app bootstrap.Application, router *gin.RouterGroup) {
	router.POST("/send", sendTxHandler(app, env))
}