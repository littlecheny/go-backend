package route

import (
	"time"
	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/littlecheny/go-backend/mongo"
	"github.com/littlecheny/go-backend/repository"
	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/domain"
	"github.com/littlecheny/go-backend/api/controller"
	"github.com/littlecheny/go-backend/usecase"
)

func NewSignupRouter(env *bootstrap.Env, db mongo.Database, timeout time.Duration, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(ur, timeout),
		Env: env,
	}
	
	group.POST("/signup", sc.Signup)
}