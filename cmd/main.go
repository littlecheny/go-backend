package main

import(
	"github.com/gin-gonic/gin"
	route "github.com/littlecheny/go-backend/route"
	"github.com/littlecheny/go-backend/bootstrap"
)

func main(){
	app := bootstrap.App()

	env := app.Env



	r := gin.Defalut()

	route.Setup(env,db,gin,timeout)

	r.Run(env.ServerAddress)
}