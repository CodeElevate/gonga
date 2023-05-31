package bootstrap

import (
	middlewares "gonga/app/Http/Middlewares"
	"gonga/database"
	_ "gonga/docs"
	"gonga/packages"
	"gonga/routes"
	"net/http"

	"github.com/pterm/pterm"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

// Application represents the Golang application instance.
type Application struct {
	Router *packages.MyRouter
	DB     *gorm.DB
}

// NewApplication creates a new instance of the Golang application.
func NewApplication() *Application {
	app := &Application{
		Router: packages.NewRouter(),
	}
	return app
}

// RegisterRoutes registers the application's routes.
func (app *Application) RegisterApiRoutes() {
	// default middlewares
	app.Router.Use(middlewares.CorsMiddleware).StrictSlash(true)
	app.Router.Use(middlewares.ThrottleMiddleware).StrictSlash(true)
	app.Router.Use(middlewares.LogMiddleware).StrictSlash(true)
	
	// Serve swagger UI
	app.Router.PathPrefix("/docs/").Handler(httpSwagger.Handler())

	routes.RegisterApiRoutes(app.Router, app.DB)
}

// ConnectDatabase connects to database.
func (app *Application) ConnectDatabase() error {
	var err error
	app.DB, err = database.Connect()
	if err != nil {
		return err
	}
	return nil
}

// Run starts the Golang application.
func (app *Application) Run() error {
	pterm.Info.Println("Server started on [http://localhost:8080]")
	return http.ListenAndServe(":8080", app.Router)
}
