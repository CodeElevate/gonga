package routes

import (
	auth "gonga/app/Http/Controllers/Auth"
	middlewares "gonga/app/Http/Middlewares"
	"gonga/packages"

	"gorm.io/gorm"
)

//	@title			Gonga Api
//	@description	A social media api.
//	@version		1.0
//	@contact.name	API Support
//	@host			localhost:8000
//	@BasePath		/
func RegisterAuthRoutes(router *packages.MyRouter, db *gorm.DB) {

	LoginController := auth.LoginController{DB: db}
	RegisterController := auth.RegisterController{DB: db}
	NewPasswordController := auth.NewPasswordController{DB: db}
	PasswordResetLinkController:=auth.PasswordResetLinkController{DB: db}

	// Login API endpoint handlers
	router.Post("/login", LoginController.Create)

	// Logout API endpoint handlers
	router.Post("/logout", LoginController.Delete, middlewares.AuthMiddleware)

	// Register API endpoint handlers
	router.Post("/register", RegisterController.Create)

	// Password reset API endpoint handlers
	router.Post("/forgot-password", PasswordResetLinkController.Create).Name("password.email")
	router.Post("/reset-password", NewPasswordController.Create).Name("password.update")
	
}
