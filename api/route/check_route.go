package route

import (
	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/bootstrap"
)

func NewCheckRouter(env *bootstrap.Env, router *gin.RouterGroup) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
}
