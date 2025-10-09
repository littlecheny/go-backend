package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/littlecheny/go-backend/mongo"
)

func Setup(env *bootstrap.Env, db mongo.Database, gin *gin.Engine, timeout time.Duration, app bootstrap.Application) {
	publicRouter := gin.Group("")

	NewSignupRouter(env, db, timeout, publicRouter)
	NewLoginRouter(env, db, timeout, publicRouter)
	NewCheckRouter(env, publicRouter)
	NewGetstatusRouter(env, publicRouter)

	chainRouter := gin.Group("/chain")
	NewChainBlockRouter(env, app, chainRouter)

	walletRouter := gin.Group("/wallet")
	NewWalletBalanceRouter(env, app, walletRouter)
}
