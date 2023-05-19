package routes

import (
	controllers "gonga/app/Http/Controllers"
	middlewares "gonga/app/Http/Middlewares"
	"gonga/packages"

	"gorm.io/gorm"
)

// @title Gonga Api
// @description A social media api.
// @version 1.0
// @contact.name API Support
// @host localhost:8000
// @BasePath /
func RegisterApiRoutes(router *packages.MyRouter, db *gorm.DB) {

	UserController := controllers.UserController{DB: db}
	SearchController := controllers.SearchController{DB: db}
	PostController := controllers.PostController{DB: db}
	NotificationController := controllers.NotificationController{DB: db}
	FollowController := controllers.FollowController{DB: db}
	MediaController := controllers.MediaController{DB: db}

	router.Post("/upload", MediaController.Upload, middlewares.AuthMiddleware)
	// User API endpoint handlers
	router.Get("/users", UserController.Index)
	router.Get("/users/{username}", UserController.Show)
	router.Put("/users/{username}", UserController.Update, middlewares.AuthMiddleware)
	router.Delete("/users/{id}", UserController.Delete, middlewares.AuthMiddleware)

	// Post API endpoint handlers
	router.Get("/posts", PostController.Index)
	router.Post("/posts", PostController.Create, middlewares.AuthMiddleware) //, middlewares.AuthMiddleware
	router.Get("/posts/{id}", PostController.Show)
	router.Put("/posts/{id}", PostController.Update, middlewares.AuthMiddleware)
	router.Put("/posts/{id}/title", PostController.UpdateTitle, middlewares.AuthMiddleware)
	router.Put("/posts/{id}/body", PostController.UpdateBody, middlewares.AuthMiddleware)
	router.Put("/posts/{id}/medias", PostController.UpdateMedia, middlewares.AuthMiddleware)
	router.Put("/posts/{id}/hashtags", PostController.UpdateHashtag, middlewares.AuthMiddleware)

	// router.Put("/posts/{id}/visibility", PostController.Update, middlewares.AuthMiddleware)
	// router.Put("/posts/{id}/is_promoted", PostController.Update, middlewares.AuthMiddleware)
	// router.Put("/posts/{id}/is_featured", PostController.Update, middlewares.AuthMiddleware)
	// router.Put("/posts/{id}/promotion_expiry", PostController.Update, middlewares.AuthMiddleware)
	// router.Put("/posts/{id}/featured_expiry", PostController.Update, middlewares.AuthMiddleware)

	
	router.Delete("/posts/{id}", PostController.Delete, middlewares.AuthMiddleware)
	// router.Get("/posts/{id}/comments", PostController.Comments, middlewares.AuthMiddleware)
	// router.Post("/posts/{id}/comments", PostController.CreateComment, middlewares.AuthMiddleware)
	// router.Post("/posts/{id}/like", PostController.Like, middlewares.AuthMiddleware)
	// router.Post("/posts/{id}/unlike", PostController.Unlike, middlewares.AuthMiddleware)

	// Follow API endpoint handlers
	router.Post("/users/{id}/friend_requests", FollowController.Index, middlewares.AuthMiddleware)

	// Notification API endpoint handlers
	router.Get("/notifications", NotificationController.Index, middlewares.AuthMiddleware)
	router.Post("/notifications/read_all", NotificationController.ReadAll, middlewares.AuthMiddleware)
	router.Post("/notifications/{id}/read", NotificationController.Update, middlewares.AuthMiddleware)

	// Search API endpoint handlers
	router.Get("/search", SearchController.Index)

	// ******************************
	// *    ALERT: DO NOT EDIT!     *
	// * This area is non-editable. *
	// ******************************

	// Register Auth Routes
	RegisterAuthRoutes(router, db)
	

}
