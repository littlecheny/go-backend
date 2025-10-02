package main

import(
	"time"
	"github.com/gin-gonic/gin"
	route "github.com/littlecheny/go-backend/api/route"
	"github.com/littlecheny/go-backend/bootstrap"
)

func main(){
	app := bootstrap.App()

	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	r := gin.Default()

	route.Setup(env, db, r, timeout)

	r.Run(env.ServerAddress)
}