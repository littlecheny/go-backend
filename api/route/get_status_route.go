package route

import (
	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/bootstrap"
)

func NewGetstatusRouter(env *bootstrap.Env, router *gin.RouterGroup) {
	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app":         env.AppEnv,
			"redisstatus": env.RedisHost,
		})
	})
}
