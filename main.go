package main

import (
	console "gonga/app/Console"
	"gonga/bootstrap"
	"log"
)

var App *bootstrap.Application

//	@title			Gonga API Documentation
//	@version		1.0
//	@description	This is the Swagger documentation for the Gonga API.
//	@host			https://gonga.up.railway.app
//	@BasePath		/
//	@contact.name	Krishan Kumar
//	@contact.email	your.email@example.com
//	@contact.url	https://www.linkedin.com/in/kkumar-gcc
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
