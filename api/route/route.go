package route

import (
	"github.com/littlecheny/go-backend/bootstrap"
)

func Setup(env *bootstrap.Env, db mongo.Database, gin *gin.Engine, timeout time.Duration){
	publicRouter := gin.Group("")

	NewSignupRouter(env, db, timeout, publicRouter)
}