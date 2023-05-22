package main

import (
	console "gonga/app/Console"
	"gonga/bootstrap"
	"log"
)

var App *bootstrap.Application


// @title Gonga API Documentation
// @version 1.0
// @description This is the Swagger documentation for the Gonga API.
// @host localhost:8080
// @BasePath /api/v1
// @contact.name Your Name
// @contact.email your.email@example.com
// @contact.url http://your-website.com
func main() {
	bootstrap.LoadEnv()

	App = bootstrap.NewApplication()

	//connect to database
	err := App.ConnectDatabase()

	if err != nil {
		log.Fatal(err)
	}

	// Register the routes
	App.RegisterApiRoutes()

	// register commands
	console.RegisterCommands(App)
}
