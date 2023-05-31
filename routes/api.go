package routes

import (
	controllers "gonga/app/Http/Controllers"
	middlewares "gonga/app/Http/Middlewares"
	"gonga/packages"

	"gorm.io/gorm"
)

// RegisterApiRoutes registers the API routes
func RegisterApiRoutes(router *packages.MyRouter, db *gorm.DB) {
	// Initialize the required controllers
	UserController := controllers.UserController{DB: db}
	SearchController := controllers.SearchController{DB: db}
	PostController := controllers.PostController{DB: db}
	NotificationController := controllers.NotificationController{DB: db}
	FollowController := controllers.FollowController{DB: db}
	MediaController := controllers.MediaController{DB: db}
	CommentController := controllers.CommentController{DB: db}
	LikeController := controllers.LikeController{DB: db}

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
	// router.Put("/posts/{id}", PostController.Update, middlewares.AuthMiddleware)
	router.Put("/posts/{id}/title", PostController.UpdateTitle, middlewares.AuthMiddleware)
	router.Put("/posts/{id}/body", PostController.UpdateBody, middlewares.AuthMiddleware)
	router.Put("/posts/{id}/medias", PostController.UpdateMedia, middlewares.AuthMiddleware)
	router.Put("/posts/{id}/hashtags", PostController.UpdateHashtag, middlewares.AuthMiddleware)
	router.Put("/posts/{id}/settings", PostController.UpdatePostSettings, middlewares.AuthMiddleware)
	router.Delete("/posts/{id}", PostController.Delete, middlewares.AuthMiddleware)

	// Comment API endpoint handlers
	router.Get("/posts/{id}/comments", CommentController.Index)
	router.Post("/posts/{id}/comments", CommentController.Create, middlewares.AuthMiddleware)
	router.Get("/comments/{id}", CommentController.Show)
	router.Put("/comments/{id}", CommentController.Update, middlewares.AuthMiddleware)
	router.Delete("/comments/{id}", CommentController.Delete, middlewares.AuthMiddleware)

	//like API endpoint handlers
	router.Post("/likes", LikeController.Create, middlewares.AuthMiddleware)
	router.Delete("/likes/{id}", LikeController.Delete, middlewares.AuthMiddleware)

	// Follow API endpoint handlers
	router.Post("/users/follow", FollowController.Create, middlewares.AuthMiddleware)

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
