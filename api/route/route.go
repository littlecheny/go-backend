package route

import (
	"time"

	"github.com/gin-gonic/gin"
	githubSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	docs "github.com/littlecheny/go-backend/docs"
	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/littlecheny/go-backend/mongo"
)

func Setup(env *bootstrap.Env, db mongo.Database, gin *gin.Engine, timeout time.Duration, app bootstrap.Application) {
	publicRouter := gin.Group("")

	// Swagger UI 路由
	docs.SwaggerInfo.BasePath = "/"
	gin.GET("/swagger/*any", githubSwagger.WrapHandler(swaggerFiles.Handler))

	NewSignupRouter(env, db, timeout, publicRouter)
	NewLoginRouter(env, db, timeout, publicRouter)
	NewCheckRouter(env, publicRouter)
	NewGetstatusRouter(env, app, publicRouter)

	chainRouter := gin.Group("/chain")
	NewChainBlockRouter(env, app, chainRouter)

	walletRouter := gin.Group("/wallet")
	NewWalletBalanceRouter(env, app, walletRouter)

	// Tx 路由分组
	txRouter := gin.Group("/tx")
	NewTxSendRouter(env, app, txRouter)
	NewTxStatusRouter(env, app, txRouter)
	NewTxTaskRouter(env, app, txRouter)
}
