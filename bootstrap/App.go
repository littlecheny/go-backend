package bootstrap

import (
	"github.com/littlecheny/go-backend/domain"
	"github.com/littlecheny/go-backend/mongo"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Env             *Env
	Mongo           mongo.Client
	Redis           *redis.Client
	EthereumService domain.EthereumService
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Mongo = NewMongoDatabase(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Mongo)
}
