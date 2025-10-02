package bootstrap

import "github.com/littlecheny/go-backend/mongo"

type Application struct{
	Env *Env
	Mongo mongo.Client
}

func App() Application{
	app := &Application{}
	app.Env = NewEnv()
	app.Mongo = NewMongoDatabase(app.Env)
	return *app
}

func (app *Application) CloseDBConnection(){
	CloseMongoDBConnection(app.Mongo)
}