package route

import (
	"time"
)

func NewSignupRouter(env *bootstrap.Env, db mongo.Database, timeout time.Duration, group *gin.RouteGroup) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(ur, timeout),
		Env: env,
	}
	
	group.POST("/signup", sc.Signup)
}