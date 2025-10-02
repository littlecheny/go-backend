package route

import (
	"time"
	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/mongo"
)

func Setup(env *bootstrap.Env, db mongo.Database, gin *gin.Engine, timeout time.Duration){
	publicRouter := gin.Group("")

	NewSignupRouter(env, db, timeout, publicRouter)
	NewLoginRouter(env, db, timeout, publicRouter)
	NewRefreshTokenRouter(env, db, timeout, publicRouter)
}