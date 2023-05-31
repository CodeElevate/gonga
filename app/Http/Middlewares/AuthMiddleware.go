package middlewares

import (
	"context"
	"errors"
	"gonga/utils"
	"net/http"
	// "strings"
)

// AuthMiddleware is a middleware function that checks if a user is authenticated.
// If the user is not authenticated, it returns an error response with status code 401.
// If the user is authenticated, it calls the next middleware/handler in the chain.
//
// Example usage:
//
//	router := packages.NewRouter()
//	router.Put("/api/users/{id}", UserController.Update, middlewares.AuthMiddleware)
//
// This will add the AuthMiddleware to the UserHandler function when the "/api/users" endpoint is accessed with the "GET" method.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const (
			userIDKey utils.ContextKey = "userID"
		)
		// Check if user is authenticated and get the user ID
		authenticated, userID := utils.IsAuthenticate(r)
		if !authenticated {
			utils.HandleError(w, errors.New("unauthorized"), http.StatusUnauthorized)
			return
		}

		// Set the user ID in the request context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		r = r.WithContext(ctx)

		// Call the next middleware/handler
		next.ServeHTTP(w, r)
	})
}
