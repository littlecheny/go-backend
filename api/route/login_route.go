package route

import(
	"time"
	"github.com/littlecheny/go-backend/repository"
	"github.com/littlecheny/go-backend/domain"
	"github.com/littlecheny/go-backend/mongo"
	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/littlecheny/go-backend/api/controller"
	"github.com/littlecheny/go-backend/usecase"
)

func NewLoginRouter(env *bootstrap.Env, db mongo.Database, timeout time.Duration, group *gin.RouterGroup){
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	sc := controller.LoginController{
		LoginUsecase: usecase.NewLoginUsecase(ur, timeout),
		Env: env,
	}

	group.POST("/login", sc.Login)
}