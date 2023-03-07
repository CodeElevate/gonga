package middlewares

import (
	"net/http"
)

func LogMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        
        next.ServeHTTP(w, r)
    })
}

