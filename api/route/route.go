package route

import (
	"time"

	"github.com/gin-gonic/gin"
	githubSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	docs "github.com/littlecheny/go-backend/docs"
	"github.com/littlecheny/go-backend/api/middleware"
	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/littlecheny/go-backend/mongo"
)

func Setup(env *bootstrap.Env, db mongo.Database, gin *gin.Engine, timeout time.Duration, app bootstrap.Application) {
	// 全局中间件
	gin.Use(middleware.CorsMiddleware())
	gin.Use(middleware.ErrorHandlerMiddleware())
	
	// 创建限流器：每分钟100个请求
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	gin.Use(rateLimiter.RateLimitMiddleware())

	// 公共路由（无需认证）
	publicRouter := gin.Group("")

	// Swagger UI 路由
	docs.SwaggerInfo.BasePath = "/"
	gin.GET("/swagger/*any", githubSwagger.WrapHandler(swaggerFiles.Handler))

	NewSignupRouter(env, db, timeout, publicRouter)
	NewLoginRouter(env, db, timeout, publicRouter)
	NewCheckRouter(env, app, publicRouter)
	NewGetstatusRouter(env, app, publicRouter)

	// 需要认证的路由
	protectedRouter := gin.Group("")
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))

	chainRouter := protectedRouter.Group("/chain")
	NewChainBlockRouter(env, app, chainRouter)

	walletRouter := protectedRouter.Group("/wallet")
	NewWalletBalanceRouter(env, app, walletRouter)

	// Tx 路由分组（需要认证）
	txRouter := protectedRouter.Group("/tx")
	NewTxSendRouter(env, app, txRouter)
	NewTxStatusRouter(env, app, txRouter)
	NewTxTaskRouter(env, app, txRouter)
}
